package main

import (
	"bytes"
	"fmt"
	"strings"
)

func genRPCHandlerInterface(iface *MethodsCollection) string {
	buf := new(bytes.Buffer)
	iface.ForEachMethod(func(m *Method) error {
		fmt.Fprintf(buf, "%s(c *gin.Context)\n", m.Name)
		return nil
	})
	return buf.String()
}

func genRPCHandlerImplementation(recvName, featurePrefix string, iface *MethodsCollection) string {
	buf := new(bytes.Buffer)
	iface.ForEachMethod(func(m *Method) error {
		fmt.Fprintf(buf, "%s\n", reqModelJSON(m))
		fmt.Fprintf(buf, "%s\n", respModelJSON(featurePrefix, m))
		fmt.Fprintf(buf, "%s\n", rpcHandlerMethod(recvName, featurePrefix, m))
		return nil
	})
	return buf.String()
}

func genServiceClientImplementation(recvName, featurePrefix string, iface *MethodsCollection) string {
	buf := new(bytes.Buffer)
	iface.ForEachMethod(func(m *Method) error {
		fmt.Fprintf(buf, "%s\n", serviceClientMethod(recvName, m))
		return nil
	})
	return buf.String()
}

func rpcHandlerMethod(recvName, featurePrefix string, m *Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "func (h *%s) %s(c *gin.Context) {\n", recvName, m.Name)

	// logging and metrics
	fmt.Fprintln(buf, `metrics.ReportFuncCall(h.tags)
	fnLog := log.WithFields(logging.WithFn(h.fields))
	statsFn := metrics.ReportFuncTiming(h.tags)
	defer statsFn()
	`)

	// request decoding
	fmt.Fprintf(buf, `var req %sRequest
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&req)
	if err != nil {
		fnLog.Warnln(err)
		bugsnag.Notify(err)
		c.JSON(http.StatusBadRequest, &%sResponse{
			%sErrorResponse: new%sErrorResponse(err),
		})
		return
	}`, m.Name, m.Name, featurePrefix, featurePrefix)
	fmt.Fprintln(buf, "")

	fmt.Fprintf(buf, "var resp %sResponse\n", m.Name)
	fmt.Fprintf(buf, "%s\n", funcCallMapping(m))

	// error handling
	fmt.Fprintf(buf, `if err != nil {
		fnLog.Warnln(err)
		bugsnag.Notify(err)
		c.JSON(http.StatusBadRequest, &%sResponse{
			%sErrorResponse: new%sErrorResponse(err),
		})
		return
	}`, m.Name, featurePrefix, featurePrefix)
	fmt.Fprintln(buf, "")

	// response encoding
	fmt.Fprintf(buf, `c.JSON(http.StatusOK, &resp)
	return
	`)

	fmt.Fprintln(buf, "}")
	return buf.String()
}

func serviceClientMethod(recvName string, m *Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "func (s *%s) %s {\n", recvName, funcSpec(m, true))

	// logging and metrics
	fmt.Fprintf(buf, `metrics.ReportFuncCall(s.tags)
	statsFn := metrics.ReportFuncTiming(s.tags)
	defer statsFn()`)
	fmt.Fprintln(buf, "\n")

	// request mapping
	fmt.Fprintf(buf, `req := %s`, reqFieldsMap(m))
	fmt.Fprintf(buf, "var resp *%sResponse\n", m.Name)

	// request
	if !hasErr(m.Res) {
		fmt.Fprintf(buf, `respBody, err := s.do(s.newReq("POST", "%s", req))`, m.Name)
	} else {
		fmt.Fprintf(buf, "var respBody []byte\n")
		fmt.Fprintf(buf, `respBody, err = s.do(s.newReq("POST", "%s", req))`, m.Name)
	}
	fmt.Fprintln(buf, "")

	// error handling
	fmt.Fprintf(buf, `if err != nil {
		return
	} else if err = s.checkRespErr(respBody); err != nil {
		return
	} else if json.Unmarshal(respBody, &resp); err != nil {
		return
	}`)
	fmt.Fprintln(buf, "")

	// response mapping
	fmt.Fprintf(buf, "%s", respFieldsMap(m))
	fmt.Fprintln(buf, `return
	}`)
	return buf.String()
}

func reqModelJSON(m *Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "type %sRequest struct {\n", m.Name)
	for _, p := range m.Params {
		if len(p.Name) == 0 {
			continue
		}
		fmt.Fprintf(buf, "%s %s `json:\"%s,omitempty\"`\n", strings.Title(p.Name), p.Type, p.Name)
	}
	fmt.Fprintln(buf, "}")
	return buf.String()
}

func reqFieldsMap(m *Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "&%sRequest {\n", m.Name)
	for _, p := range m.Params {
		if len(p.Name) == 0 {
			continue
		}
		fmt.Fprintf(buf, "%s: %s,\n", strings.Title(p.Name), p.Name)
	}
	fmt.Fprintln(buf, "}")
	return buf.String()
}

func respModelJSON(featureName string, m *Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "type %sResponse struct {\n", m.Name)
	fmt.Fprintf(buf, "%sErrorResponse\n", featureName)
	for i, p := range m.Res {
		if p.Type == "error" {
			// errors are handled by method implementation,
			// any error value will be wrapped into embedded xxxErrorResponse.
			continue
		}
		if len(p.Name) > 0 {
			fmt.Fprintf(buf, "%s %s `json:\"%s,omitempty\"`\n", strings.Title(p.Name), p.Type, p.Name)
		} else {
			fmt.Fprintf(buf, "Ret%d %s `json:\"_ret%d,omitempty\"`\n", i, p.Type, i)
		}
	}
	fmt.Fprintln(buf, "}")
	return buf.String()
}

func respFieldsMap(m *Method) string {
	buf := new(bytes.Buffer)
	for i, r := range m.Res {
		if r.Type == "error" {
			continue
		}
		if len(r.Name) == 0 {
			fmt.Fprintf(buf, "_ret%d = resp.Ret%d\n", i, i)
		} else {
			fmt.Fprintf(buf, "%s = resp.%s\n", r.Name, strings.Title(r.Name))
		}
	}
	return buf.String()
}

func funcCallMapping(m *Method) string {
	paramList := make([]string, 0, len(m.Params))
	for _, p := range m.Params {
		if len(p.Name) == 0 {
			// try nil or gtfo
			paramList = append(paramList, "nil")
		} else {
			paramList = append(paramList, fmt.Sprintf("req.%s", strings.Title(p.Name)))
		}
	}
	paramSpec := strings.Join(paramList, ", ")
	retList := make([]string, 0, len(m.Res))
	for i, r := range m.Res {
		if r.Type == "error" {
			retList = append(retList, "err")
			continue
		}
		if len(r.Name) == 0 {
			retList = append(retList, fmt.Sprintf("resp.Ret%d", i))
			continue
		}
		retList = append(retList, "resp."+strings.Title(r.Name))
	}
	retSpec := strings.Join(retList, ", ")
	if len(m.Res) == 0 {
		return fmt.Sprintf("h.svc.%s(%s)", m.Name, paramSpec)
	}
	return fmt.Sprintf("%s = h.svc.%s(%s)", retSpec, m.Name, paramSpec)
}

func funcSpec(m *Method, forceNaming bool) string {
	var retSpec string
	switch {
	case len(m.Res) == 0:
	case len(m.Res) == 1:
		ret := m.Res[0]
		if len(ret.Name) == 0 {
			if forceNaming {
				if ret.Type == "error" {
					retSpec = "(err error)"
				} else {
					retSpec = fmt.Sprintf("(_ret0 %s)", ret.Type)
				}
			} else {
				retSpec = ret.Type
			}
		} else {
			if ret.Type == "error" && forceNaming {
				retSpec = fmt.Sprintf("(err %s)", ret.Type)
			} else {
				retSpec = fmt.Sprintf("(%s %s)", ret.Name, ret.Type)
			}
		}
	default:
		rets := make([]string, 0, len(m.Res))
		if len(m.Res[0].Name) == 0 {
			if forceNaming {
				for i, r := range m.Res {
					if r.Type == "error" {
						rets = append(rets, fmt.Sprintf("err %s", r.Type))
					} else {
						rets = append(rets, fmt.Sprintf("_ret%d %s", i, r.Type))
					}
				}
			} else {
				for _, r := range m.Res {
					rets = append(rets, r.Type)
				}
			}
		} else {
			for _, r := range m.Res {
				if r.Type == "error" && forceNaming {
					rets = append(rets, "err "+r.Type)
				} else {
					rets = append(rets, r.Name+" "+r.Type)
				}
			}
		}
		retSpec = fmt.Sprintf("(%s)", strings.Join(rets, ", "))
	}
	params := make([]string, 0, len(m.Params))
	for _, p := range m.Params {
		params = append(params, p.Name+" "+p.Type)
	}
	return fmt.Sprintf("%s(%s) %s", m.Name, strings.Join(params, ", "), retSpec)
}

func hasErr(rets []Param) bool {
	for i := range rets {
		if rets[i].Type == "error" {
			return true
		}
	}
	return false
}

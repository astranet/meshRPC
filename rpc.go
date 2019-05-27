package main

import (
	"bytes"
	"fmt"
	"strings"
)

func genRPCHandlerInterface(iface *MethodsCollection) string {
	buf := new(bytes.Buffer)
	iface.ForEachMethod(func(m *Method) error {
		fmt.Fprintf(buf, "%s(ctx *gin.Context)\n", m.Name)
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
	fmt.Fprintf(buf, "func (_handler *%s) %s(_ctx *gin.Context) {\n", recvName, m.Name)

	// logging and metrics
	fmt.Fprintf(buf, "// TODO: Report Stats + Timing\n\n")

	// request decoding
	fmt.Fprintf(buf, `var _req %sRequest
	_decoder := json.NewDecoder(_ctx.Request.Body)
	defer _ctx.Request.Body.Close()
	_err := _decoder.Decode(&_req)
	if _err != nil {
		// TODO: Report Error

		_data, _ := json.Marshal(&%sResponse{
			%sErrorResponse: %sErrorResponse{
				Error: _err.Error(),
			},
		})
		_ctx.Status(400)
		// TODO: Report Stats
		_, _  = io.Copy(_ctx.Writer, bytes.NewReader(_data))
		return
	}`, m.Name, m.Name, featurePrefix, featurePrefix)
	fmt.Fprintln(buf, "")

	fmt.Fprintf(buf, "var _resp %sResponse\n", m.Name)
	fmt.Fprintf(buf, "%s\n", funcCallMapping(m))

	// error handling
	fmt.Fprintf(buf, `if _err != nil {
		// TODO: Report Error

		_data, _ := json.Marshal(&%sResponse{
			%sErrorResponse: %sErrorResponse{
				Error: _err.Error(),
			},
		})
		_ctx.Status(400)
		// TODO: Report Stats
		_, _  = io.Copy(_ctx.Writer, bytes.NewReader(_data))
		return
	}`, m.Name, featurePrefix, featurePrefix)
	fmt.Fprintln(buf, "")

	// response encoding
	fmt.Fprintf(buf, `
		_data, _ := json.Marshal(&_resp)
		_ctx.Status(200)
		// TODO: Report Stats
		_, _  = io.Copy(_ctx.Writer, bytes.NewReader(_data))
		return
	`)

	fmt.Fprintln(buf, "}")
	return buf.String()
}

func serviceClientMethod(recvName string, m *Method) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "func (_client *%s) %s {\n", recvName, funcSpec(m, true))

	// logging and metrics
	fmt.Fprintf(buf, `// TODO: Report Stats + Timing`)
	fmt.Fprintln(buf, "\n")

	// request mapping
	fmt.Fprintf(buf, `_req := %s`, reqFieldsMap(m))
	fmt.Fprintf(buf, "var _resp *%sResponse\n", m.Name)

	// request
	if !hasErr(m.Res) {
		fmt.Fprintf(buf, `_respBody, _err := _client.do(_client.newJsonReq("POST", "%s", _req))`, m.Name)
	} else {
		fmt.Fprintf(buf, "var _respBody []byte\n")
		fmt.Fprintf(buf, `_respBody, _err = _client.do(_client.newJsonReq("POST", "%s", _req))`, m.Name)
	}
	fmt.Fprintln(buf, "")

	// error handling
	fmt.Fprintf(buf, `if _err != nil {
		return
	} else if _err = _client.checkJsonRespErr(_respBody); _err != nil {
		return
	} else if json.Unmarshal(_respBody, &_resp); _err != nil {
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
			fmt.Fprintf(buf, "_ret%d = _resp.Ret%d\n", i, i)
		} else {
			fmt.Fprintf(buf, "%s = _resp.%s\n", r.Name, strings.Title(r.Name))
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
			paramList = append(paramList, fmt.Sprintf("_req.%s", strings.Title(p.Name)))
		}
	}
	paramSpec := strings.Join(paramList, ", ")
	retList := make([]string, 0, len(m.Res))
	for i, r := range m.Res {
		if r.Type == "error" {
			retList = append(retList, "_err")
			continue
		}
		if len(r.Name) == 0 {
			retList = append(retList, fmt.Sprintf("_resp.Ret%d", i))
			continue
		}
		retList = append(retList, "_resp."+strings.Title(r.Name))
	}
	retSpec := strings.Join(retList, ", ")
	if len(m.Res) == 0 {
		return fmt.Sprintf("_handler.svc.%s(%s)", m.Name, paramSpec)
	}
	return fmt.Sprintf("%s = _handler.svc.%s(%s)", retSpec, m.Name, paramSpec)
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
					retSpec = "(_err error)"
				} else {
					retSpec = fmt.Sprintf("(_ret0 %s)", ret.Type)
				}
			} else {
				retSpec = ret.Type
			}
		} else {
			if ret.Type == "error" && forceNaming {
				retSpec = fmt.Sprintf("(_err %s)", ret.Type)
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
						rets = append(rets, fmt.Sprintf("_err %s", r.Type))
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
					rets = append(rets, "_err "+r.Type)
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

// Package logging contains advanced logging patterns and helpers.
package logging

import (
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Level(s string) log.Level {
	switch s {
	case "1", "error":
		return log.ErrorLevel
	case "2", "warn":
		return log.WarnLevel
	case "3", "info":
		return log.InfoLevel
	case "4", "debug":
		return log.DebugLevel
	default:
		return log.FatalLevel
	}
}

func WithFn(fields ...log.Fields) log.Fields {
	if len(fields) > 0 && fields[0] != nil {
		result := copyFields(fields[0])
		result["fn"] = getCallerName()
		return result
	}
	return log.Fields{
		"fn": getCallerName(),
	}
}

func WithMore(fields log.Fields, add log.Fields) log.Fields {
	fields = copyFields(fields)
	for k, v := range add {
		fields[k] = v
	}
	return fields
}

func copyFields(fields log.Fields) log.Fields {
	ff := make(log.Fields, len(fields))
	for k, v := range fields {
		ff[k] = v
	}
	return ff
}

func FnName() string {
	return getCallerName()
}

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")
	return nameParts[len(nameParts)-1]
}

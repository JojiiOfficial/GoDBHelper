package godbhelper

import (
	"strings"
)

func strHasEmpty(arg ...string) bool {
	for _, s := range arg {
		if len(strings.Trim(s, " ")) == 0 {
			return true
		}
	}
	return false
}

func parseDSNstring(arg ...string) string {
	if len(arg) > 0 {
		dsn := "?" + strings.Join(arg, "&")
		return dsn
	}
	return ""
}

func isPortValid(port uint16) bool {
	return port > 0 && port <= 65535
}

func stringArrToInterface(str []string) []interface{} {
	params := make([]interface{}, len(str))
	for i, p := range str {
		params[i] = p
	}
	return params
}

func strArrHas(arr []string, has string) bool {
	for _, a := range arr {
		if a == has {
			return true
		}
	}
	return false
}

func parsetTag(tagContent string) []string {
	if strings.Contains(tagContent, ",") {
		return strings.Split(tagContent, ",")
	}
	return []string{tagContent}
}

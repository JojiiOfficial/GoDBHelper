package godbhelper

import (
	"fmt"
	"strings"
)

func buildMysqlURI(username, password, host, database string, port uint16, dsnString ...string) (string, error) {
	if strHasEmpty(username, password, host) {
		return "", ErrMysqlURIMissingArg
	}
	if !isPortValid(port) {
		return "", ErrPortInvalid
	}
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", username, password, host, port, database, parseDSNstring(dsnString...))
	return uri, nil
}

func strHasEmpty(arg ...string) bool {
	for _, s := range arg {
		if len(strings.Trim(s, " ")) == 0 {
			return true
		}
	}
	return false
}

func isPortValid(port uint16) bool {
	return port > 0 && port <= 65535
}

func parseDSNstring(arg ...string) string {
	if len(arg) > 0 {
		dsn := "?" + strings.Join(arg, "&")
		return dsn
	}
	return ""
}

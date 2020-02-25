package godbhelper

import (
	"fmt"
)

//BuildMySQLURI creates connection string for mysql
func BuildMySQLURI(username, password, host, database string, port uint16, dsnString ...string) (string, error) {
	if strHasEmpty(username, password, host) {
		return "", ErrMysqlURIMissingArg
	}

	if !isPortValid(port) {
		return "", ErrPortInvalid
	}

	return fmt.Sprintf(MysqlURIFormat, username, password, host, port, database, parseDSNstring(dsnString...)), nil
}

//BuildPostgresString creates connection string for postgres
func BuildPostgresString(username, password, host, database string, port uint16, dsnString ...string) (string, error) {
	if strHasEmpty(username, password, host) {
		return "", ErrPostgresURIMissingArg
	}

	if !isPortValid(port) {
		return "", ErrPortInvalid
	}

	return fmt.Sprintf(PostgresURIFormat, username, password, host, port, database, parsePostgresString(dsnString...)), nil
}

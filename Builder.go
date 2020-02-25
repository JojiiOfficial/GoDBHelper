package godbhelper

import (
	"fmt"
)

//BuildDSN creates connection string for mysql
func BuildDSN(dbkind dbsys, username, password, host, database string, port uint16, dsnString ...string) (string, error) {
	//Check for empty values
	if strHasEmpty(username, password, host) {
		return "", ErrMysqlURIMissingArg
	}

	//Check port
	if !isPortValid(port) {
		return "", ErrPortInvalid
	}

	//Create URI

	if dbkind == Mysql {
		return fmt.Sprintf(MysqlURIFormat, username, password, host, port, database, parseDSNstring(dsnString...)), nil
	}

	if dbkind == Postgres {
		return fmt.Sprintf(PostgresURIFormat, username, password, host, port, database, parsePostgresString(dsnString...)), nil
	}

	return "", ErrInvalidDatabase
}

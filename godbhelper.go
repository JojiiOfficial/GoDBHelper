package godbhelper

import (
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
)

//DBhelper the dbhelper object
type DBhelper struct {
	dbKind dbsys
	DB     *sqlx.DB
	IsOpen bool
}

//NewDbHelper the DBhelper constructor
func NewDbHelper(dbKind dbsys) *DBhelper {
	return &DBhelper{
		dbKind: dbKind,
	}
}

//Open db
//Sqlite 			- Open(filename)
//SqliteEncrypted	- Open(filename, key)
//Mysql  			- Open(username, password, address, port, database)
func (dbhelper *DBhelper) Open(params ...string) (*DBhelper, error) {
	switch dbhelper.dbKind {
	case Sqlite:
		{
			var dsnFlags string
			if len(params) > 1 {
				dsnFlags = parseDSNstring(params[1:]...)
			}
			dsn := "file:" + params[0] + dsnFlags
			db, err := sqlx.Open("sqlite3", dsn)
			if err != nil {
				return dbhelper, err
			}
			dbhelper.DB = db
			dbhelper.IsOpen = true
		}
	case SqliteEncrypted:
		{
			if len(params) < 2 {
				return dbhelper, ErrSqliteEncryptMissingArg
			}
			params[1] = "_crypto_key=" + params[1]
			dsnFlags := parseDSNstring(params[1:]...)
			dsn := "file:" + params[0] + dsnFlags
			db, err := sqlx.Open("sqlite3", dsn)
			if err != nil {
				return dbhelper, err
			}
			dbhelper.DB = db
			dbhelper.IsOpen = true
		}
	case Mysql:
		{
			if len(params) < 4 {
				return dbhelper, ErrMysqlURIMissingArg
			}
			dbname := ""
			if len(params) > 4 {
				dbname = params[4]
			}
			port, err := strconv.ParseUint(params[3], 10, 16)
			if err != nil {
				return dbhelper, err
			}
			uri, err := buildMysqlURI(params[0], params[1], params[2], dbname, (uint16)(port))
			if err != nil {
				return dbhelper, err
			}
			db, err := sqlx.Open("mysql", uri)
			if err != nil {
				return dbhelper, err
			}
			dbhelper.DB = db
			dbhelper.IsOpen = true
		}
	default:
		{
			return dbhelper, ErrDBNotSupported
		}
	}
	return dbhelper, nil
}

//QueryRow runs statement and fills a with db data
func (dbhelper *DBhelper) QueryRow(a interface{}, query string, args ...interface{}) error {
	return dbhelper.DB.Get(a, query, args...)
}

//QueryRows like QueryRow but for multiple rows
func (dbhelper *DBhelper) QueryRows(a interface{}, query string, args ...interface{}) error {
	return dbhelper.DB.Select(a, query, args...)

}

//Exec executes command in DB
func (dbhelper *DBhelper) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dbhelper.DB.Exec(query, args...)
}

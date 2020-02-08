package godbhelper

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//DBhelper the dbhelper object
type DBhelper struct {
	dbKind dbsys
	DB     *sqlx.DB
}

//NewDbHelper the DBhelper constructor
func NewDbHelper(dbKind dbsys) *DBhelper {
	return &DBhelper{
		dbKind: dbKind,
	}
}

//Open opens the db
//Sqlite - Open(filename)
func (dbhelper *DBhelper) Open(params ...string) error {
	switch dbhelper.dbKind {
	case Sqlite:
		{
			if dbhelper.dbKind == Sqlite {
				db, err := sqlx.Open("sqlite3", params[0])
				if err != nil {
					return err
				}
				dbhelper.DB = db
			}
		}
	default:
		{
			return ErrDBNotSupported
		}
	}
	return nil
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

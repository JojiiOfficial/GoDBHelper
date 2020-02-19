package godbhelper

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
)

//DBhelperOptions options for DBhelper
type DBhelperOptions struct {
	Debug             bool
	StopUpdateOnError bool
	StoreVersionInDB  bool
	UseColors         bool
}

//DBhelper the dbhelper object
type DBhelper struct {
	//used database system
	dbKind dbsys

	//Versions for upgrading
	//CurrentVersion the version currently running
	CurrentVersion float32
	//AvailableVersion the version which is newly added
	AvailableVersion float32

	//DBhelper data
	DB          *sqlx.DB
	Options     DBhelperOptions
	IsOpen      bool
	QueryChains []QueryChain `json:"chains"`

	ErrHookFunc    ErrHookFunc
	ErrHookOptions *ErrHookOptions

	NextErrHookFunc   ErrHookFunc
	NextErrHookOption *ErrHookOptions
	NextLogPrefix     *string
}

//NewDBHelper the DBhelper constructor NewDBHelper(database, debug, stopUpdateOnError, storeVersionInDB, useColors)
func NewDBHelper(dbKind dbsys, bv ...bool) *DBhelper {
	options := DBhelperOptions{
		StoreVersionInDB: true,
		UseColors:        true,
	}

	for i, v := range bv {
		switch i {
		case 0:
			options.Debug = v
		case 1:
			options.StopUpdateOnError = v
		case 2:
			options.StoreVersionInDB = v
		case 3:
			options.UseColors = v
		}
	}

	dbhelper := DBhelper{
		dbKind:  dbKind,
		Options: options,
	}

	return &dbhelper
}

func (dbhelper *DBhelper) initDBVersion() error {
	dbhelper.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version FLOAT)", TableDBVersion))

	var c int
	dbhelper.QueryRow(&c, "SELECT COUNT(*) FROM "+TableDBVersion)

	if c == 0 {
		dbhelper.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES (?)", TableDBVersion), -1.0)
	} else if c > 1 {
		return ErrVersionStoreTooManyVersions
	}

	//Load version into dbhelper
	dbhelper.QueryRow(&dbhelper.CurrentVersion, "SELECT version FROM "+TableDBVersion)
	return nil
}

func (dbhelper *DBhelper) saveVersion(version float32) {
	if dbhelper.Options.StoreVersionInDB {
		dbhelper.Exec("DELETE FROM " + TableDBVersion)
		dbhelper.Exec(fmt.Sprintf("INSERT INTO %s (version) VALUES (?)", TableDBVersion), version)
	}
	dbhelper.CurrentVersion = version
}

//SetErrHook sets the error hook function
func (dbhelper *DBhelper) SetErrHook(hook ErrHookFunc, options ...ErrHookOptions) {
	if len(options) > 0 {
		dbhelper.ErrHookOptions = &options[0]
	}

	dbhelper.ErrHookFunc = hook
}

//WithHook adds next log prefix
func (dbhelper *DBhelper) WithHook(hook ErrHookFunc, options ...ErrHookOptions) *DBhelper {
	if len(options) > 0 {
		dbhelper.NextErrHookOption = &options[0]
	}
	dbhelper.NextErrHookFunc = hook
	return dbhelper
}

//WithMessage adds next log prefix
func (dbhelper *DBhelper) WithMessage(s string) *DBhelper {
	dbhelper.NextLogPrefix = &s
	return dbhelper
}

//Open db
//Sqlite 			- Open(filename)
//SqliteEncrypted	- Open(filename, key)
//Mysql  			- Open(username, password, address, port, database)
func (dbhelper *DBhelper) Open(params ...string) (*DBhelper, error) {
	dbhelper.checkColors()
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
			var dbname string
			if len(params) > 4 {
				dbname = params[4]
			}
			port, err := strconv.ParseUint(params[3], 10, 16)
			if err != nil {
				return dbhelper, err
			}
			uri, err := buildMySQLURI(params[0], params[1], params[2], dbname, (uint16)(port))
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

	if dbhelper.Options.StoreVersionInDB {
		dbhelper.initDBVersion()
	} else if dbhelper.Options.Debug {
		fmt.Println("Note: No DBVersion was restored!")
	}
	return dbhelper, nil
}

//QueryRow runs statement and fills a with db data
func (dbhelper *DBhelper) QueryRow(a interface{}, query string, args ...interface{}) error {
	err := dbhelper.DB.Get(a, query, args...)
	return dbhelper.handleErrHook(err, query)
}

//QueryRowf like QueryRow but formatted
func (dbhelper *DBhelper) QueryRowf(a interface{}, query string, queryArgs []string, args ...interface{}) error {
	return dbhelper.DB.Get(a, fmt.Sprintf(query, stringArrToInterface(queryArgs)...), args...)
}

//QueryRows like QueryRow but for multiple rows
func (dbhelper *DBhelper) QueryRows(a interface{}, query string, args ...interface{}) error {
	return dbhelper.handleErrHook(dbhelper.DB.Select(a, query, args...), query)
}

//QueryRowsf like QueryRows but formatted
func (dbhelper *DBhelper) QueryRowsf(a interface{}, query string, queryArgs []string, args ...interface{}) error {
	return dbhelper.handleErrHook(dbhelper.DB.Select(a, fmt.Sprintf(query, stringArrToInterface(queryArgs)...), args...), query)
}

//Exec executes command in DB
func (dbhelper *DBhelper) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := dbhelper.DB.Exec(query, args...)
	err = dbhelper.handleErrHook(err, query)
	return res, err
}

//Execf executes a formatted query in DB
func (dbhelper *DBhelper) Execf(queryFormat string, formatParams []string, args ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf(queryFormat, stringArrToInterface(formatParams)...)
	res, err := dbhelper.DB.Exec(query, args...)
	err = dbhelper.handleErrHook(err, query)
	return res, err
}

//AddQueryChain adds a queryChain
func (dbhelper *DBhelper) AddQueryChain(chain QueryChain) *DBhelper {
	dbhelper.QueryChains = append(dbhelper.QueryChains, chain)
	return dbhelper
}

//RunUpdate updates new sql queries
//RunUpdate(fullUpdate, dropAllTables bool)
func (dbhelper *DBhelper) RunUpdate(options ...bool) error {
	if !dbhelper.Options.StoreVersionInDB {
		return ErrCantStoreVersionInDB
	}
	dbhelper.checkColors()

	var fullUpdate, dropAllTables bool
	for i, v := range options {
		switch i {
		case 0:
			fullUpdate = v
		case 1:
			dropAllTables = v
		}
	}

	if dbhelper.Options.Debug {
		var add string
		if fullUpdate {
			add = "full"
		}
		fmt.Printf("Updating database %s\n", add)
	}

	if fullUpdate {
		dbhelper.CurrentVersion = 0
	}

	if dropAllTables {
		//TODO
	}

	var c int
	noError := true
	newVersion := dbhelper.CurrentVersion

	if dbhelper.Options.Debug {
		fmt.Println()
	}

	sort.SliceStable(dbhelper.QueryChains, func(i, j int) bool {
		return dbhelper.QueryChains[i].Order < dbhelper.QueryChains[j].Order
	})
	for _, chain := range dbhelper.QueryChains {
		if dbhelper.Options.Debug {
			color.New(color.Underline).Println("chain:", chain.Name)
		}
		sort.SliceStable(chain.Queries, func(i, j int) bool {
			return chain.Queries[i].VersionAdded < chain.Queries[j].VersionAdded
		})
		o := 0
		for _, query := range chain.Queries {
			if len(query.QueryString)+len(query.FqueryString) == 0 {
				continue
			}
			if query.VersionAdded > dbhelper.CurrentVersion {
				if dbhelper.Options.Debug {
					q := fmt.Sprintf(query.FqueryString, stringArrToInterface(query.Fparams)...)
					if len(query.FqueryString) == 0 {
						q = query.QueryString
					}
					fmt.Print("v.", query.VersionAdded, ":\t\"", q, "\"", query.Params)
				}

				var err error
				if len(query.FqueryString) > 0 {
					_, err = dbhelper.Execf(query.FqueryString, query.Fparams, stringArrToInterface(query.Params)...)
				} else {
					_, err = dbhelper.Exec(query.QueryString, stringArrToInterface(query.Params)...)
				}
				if err != nil {
					fmt.Printf(" -> %s\n", color.New(color.FgRed).SprintFunc()(" ERROR: "+err.Error()))
					if dbhelper.Options.StopUpdateOnError {
						return err
					}
					noError = false
				}
				if query.VersionAdded > newVersion {
					newVersion = query.VersionAdded
				}
				if dbhelper.Options.Debug && err == nil {
					fmt.Printf(" -> %s\n", color.New(color.FgGreen).SprintFunc()("success"))
					o++
				}
				c++
			}
		}
		if dbhelper.Options.Debug && o > 0 {
			fmt.Println()
		}
	}
	if dbhelper.Options.Debug {
		msg := "Updated %d Database queries with errors\n"
		if noError {
			msg = "Successfully updated %d Database queries\n"
		}
		fmt.Printf(msg, c)
	}
	dbhelper.saveVersion(newVersion)
	return nil
}

func (dbhelper *DBhelper) checkColors() {
	color.NoColor = !dbhelper.Options.UseColors
}

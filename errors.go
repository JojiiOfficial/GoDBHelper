package godbhelper

import "errors"

var (
	//ErrDBNotSupported error if database is not supported
	ErrDBNotSupported = errors.New("Database not supported")
	//ErrMysqlURIMissingArg error if Open() mysqldb and missing an arg
	ErrMysqlURIMissingArg = errors.New("MYSQL missing argument. Use Open(username, password, address, port, database)")
	//ErrPortInvalid if given port is invalid
	ErrPortInvalid = errors.New("Port invalid. Port must be <= 65535 and > 0")
	//ErrSqliteEncryptMissingArg errer if Open() SqliteEncrypt and missing argument
	ErrSqliteEncryptMissingArg = errors.New("SqliteEncrypt missing argument. Use Open(file, key)")
	//ErrVersionStoreTooManyVersions if VersionStore contains more than one version
	ErrVersionStoreTooManyVersions = errors.New("Too many versions stored in VersionStore")
	//ErrCantStoreVersionInDB err if runnig update and StoreVersionInDB=false
	ErrCantStoreVersionInDB = errors.New("Can't store Version in Database. Set StoreVersionInDB=true")

	//QueryBuilder errors

	//ErrNoStruct if the given data is no struct
	ErrNoStruct = errors.New("Data must be a struct")
)

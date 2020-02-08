package godbhelper

type dbsys int

const (
	//Sqlite  sqlite db
	Sqlite dbsys = iota
	//SqliteEncrypted Sqlite encrypted
	SqliteEncrypted
	//Mysql mysql db
	Mysql
	//Postgres postgres db
	Postgres
)

const (
	//TableDBVersion tablename for db version store
	TableDBVersion = "DBVersion"
)

const (
	//MysqlURIFormat formats mysql uri
	MysqlURIFormat = "%s:%s@tcp(%s:%d)/%s%s"
)

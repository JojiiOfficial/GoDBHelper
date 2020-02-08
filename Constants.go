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

package godbhelper

type dbsys int

const (
	//Sqlite  sqlite db
	Sqlite dbsys = iota
	//Mysql mysql db
	Mysql
	//Postgres postgres db
	Postgres
)

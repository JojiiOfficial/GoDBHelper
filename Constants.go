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
	//TableDBVersion tableName for db version store
	TableDBVersion = "DBVersion"
)

const (
	//MysqlURIFormat formats mysql uri
	MysqlURIFormat = "%s:%s@tcp(%s:%d)/%s%s"
	//PostgresURIFormat formats mysql uri
	PostgresURIFormat = "user='%s' password='%s' host='%s' port=%d dbname='%s' %s"
)

//Tags
const (
	//OrmTag orm-tag
	OrmTag = "orm"
	//DBTag db-tag
	DBTag = "db"
)

//Tag values
const (
	TagPrimaryKey    = "pk"
	TagAutoincrement = "ai"
	TagNotNull       = "nn"
)

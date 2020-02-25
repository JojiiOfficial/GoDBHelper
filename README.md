# GoDBHelper
A simple and lightweight orm library for golang

# Features

- Database [versioning](https://github.com/JojiiOfficial/GoDBHelper#versioning)/migrating
- Executing prepared/named/normal statements easily with formatting strings (queries)
- Easily switching between Databases (see [Driver](https://github.com/JojiiOfficial/GoDBHelper#driver))
- All [sqlx](https://github.com/jmoiron/sqlx) functions

### Driver
- [Sqlite3](https://github.com/mattn/go-sqlite3)
- [Sqlite3Encrypt](https://github.com/CovenantSQL/go-sqlite3-encrypt)
- [MySQL](https://github.com/go-sql-driver/mysql)
- [Postgres](https://github.com/lib/pq) (not completely supported yet)



# Usage
Use one of the following imports matching the driver you want to use.<br>
Sqlite: `github.com/mattn/go-sqlite3`<br>
Sqlite encrypt: `github.com/CovenantSQL/go-sqlite3-encrypt`<br>
MySQL: `github.com/go-sql-driver/mysql`<br>
PostgreSQL: `github.com/lib/pq`<br>

# Example

### Connections
```go
package main

import (
	"fmt"
	dbhelper "github.com/JojiiOfficial/GoDBHelper/"

	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
	//_ "github.com/lib/pq"
	//_ "github.com/CovenantSQL/go-sqlite3-encrypt"
)

type testUser struct {
	ID       int    `db:"id" orm:"pk,ai"`
	Username string `db:"username"`
	Pass     string `db:"password"`
}

func sqliteExample() {
	db := connectToSqlite()
	if db == nil {
		return
	}
	defer db.DB.Close()

	db.Exec("CREATE TABLE user (id int, username text, password text)")
	db.Exec("INSERT INTO user (id, username, password) VALUES (1,'will', 'iamsafe')")

	var user testUser
	db.QueryRow(&user, "SELECT * FROM user")
	fmt.Println(user.ID, ":", user.Username, user.Pass)
}


//TestStruct an example for MySQL
type TestStruct struct {
	Pkid      uint32    `db:"pk_id" orm:"pk,ai,nn"`
	Name      string    `db:"name" orm:"nn"`
	Age       uint8     `db:"age" orm:"nn" default:"1"`
	Email     string    `db:"email" orm:"nn"`
	CreatedAt time.Time `db:"createdAt" orm:"nn" default:"now()"`
}

//An example using mysql as database
func mysqlExample(){
	db := connectToMysql()
	if db == nil {
		return
	}
	defer db.DB.Close()
	
	//Create a Table from a struct. (CreateOption is optional)
	err = db.CreateTable(TestStruct{}, &godbhelper.CreateOption{
		//Create table if not exists
		IfNotExists: true,
		//Use a different name for the table than 'TestStruct'
		TableName: "TestDB",
	})
	
	s1 := TestStruct{
		Email: "email@test.com",
		Name:  "goDbHelper",
	}

	//Insert s1 into the Database. If you want to automatically set the PKid field, you have to pass the address of s1!
	resultSet, err = db.Insert(&s1, &godbhelper.InsertOption{
		//Ignore 'age' to let the DB insert the default value (otherwise it would be 0)
		IgnoreFields: []string{"age"},
		//Automatically fill the PKid field in s1. Only works if the 'orm'-Tag contains 'pk' and 'ai' and the reference to s1 is passed
		SetPK:        true,
	})
	_ = resultSet

	//Load the new entry in s2. Note that you have to set parseTime=True to read 'createdAt' in a time.Time struct
	var s2 TestStruct
	err = db.QueryRow(&s2, "SELECT * FROM TestStruct WHERE pk_id=?", s1.Pkid)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s2)
	}
}

//Connect to Mysql and return a DBhelper pointer
func connectToMysql() *dbhelper.DBhelper {
	user := "dbUser"
	pass := "pleaseMakeItSafe"
	host := "localhost"
	port := "3306"
	database := "test"
	db, err := dbhelper.NewDBHelper(dbhelper.Mysql).Open(user, pass, host, port, database, "parseTime=True")
	if err != nil {
		fmt.Fatal(err.Error())
		return nil
	}
	return db
}

func connectToSqlite() *dbhelper.DBhelper {
	db, err := dbhelper.NewDBHelper(dbhelper.Sqlite).Open("test.db")
	if err != nil {
		fmt.Fatal(err.Error())
		return nil
	}
	return db
}

func connectToSqliteEncrypt() *dbhelper.DBhelper {
	db, err := dbhelper.NewDBHelper(dbhelper.SqliteEncrypted).Open("test.db", "passKEY")
	if err != nil {
		fmt.Fatal(err.Error())
		return nil
	}
	return db
}

```
### Migrating
The following codesnippet demonstrates, how you can integrate database migration to your applications<br>

```go
//db is an instance of dbhelper.DBhelper

//load sql queries from .sql file
//Queries loaded from this function (LoadQueries) are always version 0. The last argument ('0') specifies the order of the chains.
db.LoadQueries("chain1", "./test.sql", 0)

//Add sql queries manually
//The order specifies the execution order of the queries. So in this case, chain1 would be loaded before chain2
db.AddQueryChain(dbhelper.QueryChain{
	Order: 1,
	Name: "chain2",
	Queries: []dbhelper.SQLQuery{
		dbhelper.SQLQuery{
			VersionAdded: 0,
			QueryString:  "CREATE TABLE user (id int, username text, password text)",
		},
		dbhelper.SQLQuery{
			VersionAdded: 0,
			QueryString:  "INSERT INTO user (id, username, password) VALUES (?,?,?)",
			Params:       []string{"0", "admin", "lol123"},
		},
		//added in a later version (version 0.1)
		dbhelper.SQLQuery{
			VersionAdded: 0.1,
			QueryString:  "CREATE TABLE test1 (id int)",
		},
		//added in a later version (version 0.21)
		dbhelper.SQLQuery{
			VersionAdded: 0.21,
			QueryString:  "INSERT INTO test1 (id) VALUES (?),(?)",
			Params:       []string{"29", "1"},
		},
	},
})

//runs the update
err := db.RunUpdate()
if err != nil {
	fmt.Println("Err updating", err.Error())
}
```
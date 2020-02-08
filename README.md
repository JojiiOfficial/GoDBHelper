# GoDBHelper
A database helper for golang

# Features

- Database versioning/upgrading ()
- Executing prepared/named/normal statements
- All [sqlx](https://github.com/jmoiron/sqlx) functions

### Driver
- [Sqlite3](https://github.com/mattn/go-sqlite3)
- [Sqlite3Encrypt](https://github.com/CovenantSQL/go-sqlite3-encrypt)
- [MySQL](github.com/go-sql-driver/mysql)
- Postgres

# Example

### Connections
```go
package main

import (
	"fmt"
	dbhelper "github.com/JojiiOfficial/GoDBHelper/"

	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/CovenantSQL/go-sqlite3-encrypt"
)

type testUser struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Pass     string `db:"password"`
}

func main() {
	db := exampleSqlite()
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

func exampleMysql() *dbhelper.DBhelper {
	user := "dbUser"
	pass := "pleaseMakeItSafe"
	host := "localhost"
	port := "3306"
	database := "test"
	db, err := dbhelper.NewDBHelper(dbhelper.Mysql).Open(user, pass, host, port, database)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

func exampleSqlite() *dbhelper.DBhelper {
	db, err := dbhelper.NewDBHelper(dbhelper.Sqlite).Open("test.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

func exampleSqliteEncrypt() *dbhelper.DBhelper {
	db, err := dbhelper.NewDBHelper(dbhelper.SqliteEncrypted).Open("test.db", "passKEY")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

```
### Versioning
The following code snipped demonstrates, how you can easily run new queries to upgrade the client database schema.<br>
```go
//db is an instance of dbhelper.DBhelper
db.AddQueryChain(dbhelper.QueryChain{
	Name: "chain1",
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
		//added in a later version
		dbhelper.SQLQuery{
			VersionAdded: 0.1,
			QueryString:  "CREATE TABLE test1 (id int)",
		},
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
If you add some queries in a later version, the only thing you have to do is adding a SQLQuery element to this array with a new and bigger version number. Clients wich are running on a lower version number, will run this SQL queries directly on the first run after updating (eg. `git pull` or a `docker pull`).

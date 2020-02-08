# GoDBHelper
A database helper for golang

# Features

- Database versioning/upgrading
- Executing prepared/named/normal statements
- All [sqlx](https://github.com/jmoiron/sqlx) functions

### Driver
- [Sqlite3](https://github.com/mattn/go-sqlite3)
- [Sqlite3Encrypt](https://github.com/CovenantSQL/go-sqlite3-encrypt)
- [MySQL](github.com/go-sql-driver/mysql)
- Postgres

# Example

```go
package main

import (
	"fmt"
	dbhelper "godbhelper"

	//	_ "github.com/go-sql-driver/mysql"
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
	db, err := dbhelper.NewDbHelper(dbhelper.Mysql).Open(user, pass, host, port, database)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

func exampleSqlite() *dbhelper.DBhelper {
	db, err := dbhelper.NewDbHelper(dbhelper.Sqlite).Open("test.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

func exampleSqliteEncrypt() *dbhelper.DBhelper {
	db, err := dbhelper.NewDbHelper(dbhelper.SqliteEncrypted).Open("test.db", "passKEY")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

```

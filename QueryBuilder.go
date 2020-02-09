package godbhelper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

//SQLColumn a column in a table
type SQLColumn struct {
	Name string
	Type string
}

func (dbhelper *DBhelper) create(name string, data interface{}, additionalFields ...SQLColumn) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return ErrNoStruct
	}
	if len(name) == 0 {
		name = t.Name()
	}

	v := reflect.ValueOf(data)
	var sbuff, pk string

fl:
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag

		if getSQLKind(field.Kind(), dbhelper.dbKind) == "" {
			return errors.New("Kind " + field.Kind().String() + " not supported")
		}

		colName := v.Type().Field(i).Name
		colType := getSQLKind(field.Kind(), dbhelper.dbKind)

		if len(tag.Get("db")) > 0 {
			colName = tag.Get("db")
		}

		if len(tag.Get("orm")) > 0 {
			var arr []string
			if strings.Contains(tag.Get("orm"), ",") {
				arr = strings.Split(tag.Get("orm"), ",")
			} else {
				arr = append(arr, tag.Get("orm"))
			}

			if strArrHas(arr, "-") {
				continue fl
			}

			for _, tag := range arr {
				switch tag {
				case "pk":
					pk = colName
				case "ai":
					{
						colType += " AUTO_INCREMENT"
					}
				}
			}
		}

		sbuff += "`" + colName + "` " + colType + ", "
	}

	if len(pk) > 0 {
		pk = fmt.Sprintf(", PRIMARY KEY (`%s`)", pk)
	}

	query := fmt.Sprintf("CREATE TABLE `%s` (%s%s)", name, sbuff[:len(sbuff)-2], pk)
	_, err := dbhelper.Exec(query)
	if dbhelper.Options.Debug {
		fmt.Println(query)
	}
	return err
}

func getSQLKind(kind reflect.Kind, database dbsys) string {
	switch kind {
	case reflect.String:
		return "TEXT"
	case reflect.Int:
		return "INT"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.Bool:
		return boolValue(database)
	case reflect.Int8:
		return "SMALLINT"
	case reflect.Int16:
		return "MEDIUMINT"
	case reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	default:
		return ""
	}
}

func boolValue(database dbsys) string {
	switch database {
	case Sqlite, SqliteEncrypted:
		return "INT"
	case Mysql:
		return "TINYINT(1)"
	case Postgres:
		return "boolean"
	}
	return ""
}

//CreateTable creates a table for struct
//Leave name empty to use the name of the struct
func (dbhelper *DBhelper) CreateTable(name string, data interface{}, additionalFields ...SQLColumn) error {
	return dbhelper.create(name, data, additionalFields...)
}

package godbhelper

import (
	"errors"
	"fmt"
	"reflect"
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
	var sbuff string
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if getSQLKind(f.Kind(), dbhelper.dbKind) == "" {
			return errors.New("Kind " + f.Kind().String() + " not supported")
		}
		colName := v.Type().Field(i).Name
		if len(v.Type().Field(i).Tag.Get("db")) > 0 {
			colName = v.Type().Field(i).Tag.Get("db")
		}

		colType := f.Type().Name()

		sbuff += colName + " " + colType + ", "
	}
	query := fmt.Sprintf("CREATE TABLE %s (%s)", name, sbuff[:len(sbuff)-2])
	_, err := dbhelper.Exec(query)
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

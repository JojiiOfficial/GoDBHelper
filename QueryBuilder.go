package godbhelper

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
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

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag

		if getSQLKind(field.Kind(), dbhelper.dbKind) == "" {
			return errors.New("Kind " + field.Kind().String() + " not supported")
		}

		colName := v.Type().Field(i).Name
		colType := getSQLKind(field.Kind(), dbhelper.dbKind)

		//Tags
		dbTag := tag.Get(DBTag)
		ormtag := tag.Get(OrmTag)

		if len(dbTag) > 0 {
			colName = dbTag
		}

		if len(ormtag) > 0 {
			ormTagList := parsetTag(ormtag)
			if strArrHas(ormTagList, "-") {
				continue
			}

			for _, tag := range ormTagList {
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

		colName = fmt.Sprintf("`%s`", colName)
		sbuff += fmt.Sprintf("%s %s, ", colName, colType)
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

func (dbhelper *DBhelper) insert(tableName string, data interface{}) (*sql.Result, error) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return nil, ErrNoStruct
	}
	if len(tableName) == 0 {
		tableName = t.Name()
	}

	v := reflect.ValueOf(data)
	var valuesBuff, typesBuff string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag

		colName := v.Type().Field(i).Name
		colType := field.Kind()

		var cva string
		switch colType {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			cva = strconv.FormatInt(field.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			cva = strconv.FormatUint(field.Uint(), 10)
		case reflect.String:
			cva = field.String()
		case reflect.Float64:
			cva = strconv.FormatFloat(field.Float(), 'f', 5, 64)
		case reflect.Float32:
			cva = strconv.FormatFloat(field.Float(), 'f', 5, 32)
		case reflect.Bool:
			cva = strconv.FormatBool(field.Bool())
		default:
			log.Printf("Kind %s not found!\n", colType.String())
			return nil, errors.New("Kind not supported")
		}

		//Tags
		dbTag := tag.Get(DBTag)
		ormtag := tag.Get(OrmTag)

		if len(dbTag) > 0 {
			colName = dbTag
		}

		if len(ormtag) > 0 {
			ormTagList := parsetTag(ormtag)
			if strArrHas(ormTagList, "-") {
				continue
			}
			if strArrHas(ormTagList, "ai") && !strArrHas(ormTagList, "iai") {
				continue
			}
		}

		typesBuff += fmt.Sprintf("`%s`, ", colName)

		val := cva
		switch colType {
		case reflect.String:
			val = fmt.Sprintf("'%s'", val)
		}

		valuesBuff += val + ", "
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, typesBuff[:len(typesBuff)-2], valuesBuff[:len(valuesBuff)-2])
	if dbhelper.Options.Debug {
		fmt.Println(query)
	}

	result, err := dbhelper.Exec(query)
	return &result, err
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
func float64Value(databate dbsys) string {
	switch databate {
	case Sqlite, SqliteEncrypted, Mysql:
		return "DOUBLE"
	case Postgres:
		return "numeric"
	}
	return ""
}
func float32Value(databate dbsys) string {
	switch databate {
	case Sqlite, SqliteEncrypted, Mysql:
		return "FLOAT"
	case Postgres:
		return "REAL"
	}
	return ""
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
	return dbhelper.handleErrHook(dbhelper.create(name, data, additionalFields...), "creating table "+name)
}

//Insert creates a table for struct
//Leave name empty to use the name of the struct
func (dbhelper *DBhelper) Insert(data interface{}, params ...string) (*sql.Result, error) {
	var tbName string
	if len(params) > 0 {
		tbName = params[0]
	}
	res, err := dbhelper.insert(tbName, data)
	return res, dbhelper.handleErrHook(err, "inserting "+tbName)
}

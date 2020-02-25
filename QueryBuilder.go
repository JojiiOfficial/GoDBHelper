package godbhelper

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
)

//InsertOption options for inserting structs into DB
type InsertOption struct {
	TableName        string
	IgnoreFields     []string
	SetPK            bool
	FillNotSetFields bool
}

//CreateOption options for inserting structs into DB
type CreateOption struct {
	TableName   string
	IfNotExists bool
}

func (dbhelper *DBhelper) create(data interface{}, option *CreateOption) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		return ErrNoStruct
	}

	tableName := t.Name()

	if option != nil && len(option.TableName) > 0 {
		tableName = option.TableName
	}

	v := reflect.ValueOf(data)
	var sbuff, pk string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		//fmt.Println(field.Kind(), field.Type().Kind().String(), reflect.TypeOf(time.Time{}))

		//Determine column type according to the used database
		colType := getSQLKind(field.Type(), dbhelper.dbKind)
		if colType == "" {
			return errors.New("Kind " + field.Type().String() + " not supported")
		}

		tag := v.Type().Field(i).Tag
		colName := v.Type().Field(i).Name

		//Tags
		dbTag := tag.Get(DBTag)
		ormTag := tag.Get(OrmTag)

		if len(dbTag) > 0 {
			colName = dbTag
		}

		if len(ormTag) > 0 {
			ormTagList := parsetTag(ormTag)
			if strArrHas(ormTagList, TagIgnore) {
				continue
			}

			for _, tag := range ormTagList {
				switch tag {
				case TagPrimaryKey:
					pk = colName
				case TagAutoincrement:
					{
						colType += " AUTO_INCREMENT"
					}
				case TagNotNull:
					colType += " NOT NULL"
				}
			}
		}

		//Set default value if available
		defaultTag := tag.Get(DefaultTag)
		if len(defaultTag) > 0 {
			colType += " DEFAULT " + defaultTag
		}

		colName = fmt.Sprintf("`%s`", colName)
		sbuff += fmt.Sprintf("%s %s, ", colName, colType)
	}

	if len(pk) > 0 {
		pk = fmt.Sprintf(", PRIMARY KEY (`%s`)", pk)
	}

	//Add 'if not exists' to the query if required
	tadd := ""
	if option != nil && option.IfNotExists {
		tadd = "IF NOT EXISTS"
	}

	query := fmt.Sprintf("CREATE TABLE %s `%s` (%s%s)", tadd, tableName, sbuff[:len(sbuff)-2], pk)
	_, err := dbhelper.Exec(query)
	if dbhelper.Options.Debug {
		fmt.Println(query)
	}
	return err
}

func (dbhelper *DBhelper) insert(data interface{}, option *InsertOption) (*sql.Result, error) {
	t := reflect.TypeOf(data)
	isPointer := false

	//Use Elem() if data is pointer
	if t.Kind() == reflect.Ptr {
		isPointer = true
		t = t.Elem()
	}

	//Check if data (or its value) is a struct
	if t.Kind() != reflect.Struct {
		return nil, ErrNoStruct
	}

	//Use option table name if available
	var tableName string
	if option != nil && len(option.TableName) > 0 {
		tableName = option.TableName
	} else {
		tableName = t.Name()
	}

	//Use correct reflect.Value
	v := reflect.ValueOf(data)
	if isPointer {
		v = v.Elem()
	}

	//Return error if can't address input but SetPK is required
	if !v.CanSet() && option != nil && option.SetPK {
		return nil, ErrCantAddress
	}

	var pkField *reflect.Value

	//Loop fields
	var valuesBuff, typesBuff string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag

		//Tags
		dbTag := tag.Get(DBTag)
		ormtag := tag.Get(OrmTag)

		if len(ormtag) > 0 {
			ormTagList := parsetTag(ormtag)

			//Set pkField to cur fieldAddress to set the new PK
			if option != nil && option.SetPK &&
				strArrHas(ormTagList, TagAutoincrement) && strArrHas(ormTagList, TagPrimaryKey) {
				pkField = &field
			}

			if strArrHas(ormTagList, TagIgnore) || (strArrHas(ormTagList, TagAutoincrement) && !strArrHas(ormTagList, TagInsertAutoincrement)) {
				continue
			}
		}

		colName := v.Type().Field(i).Name
		colType := field.Kind()

		if len(dbTag) > 0 {
			if dbTag == TagIgnore {
				continue
			}
			colName = dbTag
		}

		//If columnName is on Ignore list, skip it
		if option != nil && len(option.IgnoreFields) > 0 && strArrHas(option.IgnoreFields, colName) {
			continue
		}

		//Get value of field as string
		cva, defaultVal, err := strValueFromReflect(field)
		if err != nil {
			return nil, err
		}

		//Skip empty fields
		if len(cva) == 0 && (option != nil && !option.FillNotSetFields) {
			continue
		}

		//Set default value if not skipping
		if len(cva) == 0 {
			cva = defaultVal
		}

		typesBuff += "`" + colName + "`, "

		if colType == reflect.String {
			cva = fmt.Sprintf("'%s'", cva)
		}

		valuesBuff += cva + ", "
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, typesBuff[:len(typesBuff)-2], valuesBuff[:len(valuesBuff)-2])
	if dbhelper.Options.Debug {
		fmt.Println(query)
	}

	result, err := dbhelper.Exec(query)

	if pkField != nil && err == nil && result != nil {
		id, err := result.LastInsertId()
		if err != nil {
			return &result, err
		}
		if isUnsigned(pkField.Type().Kind()) {
			pkField.SetUint(uint64(id))
		} else {
			pkField.SetInt(id)
		}
	}

	return &result, err
}

func getSQLKind(kind reflect.Type, database dbsys) string {
	switch kind.Kind() {
	case reflect.String:
		return "TEXT"
	case reflect.Float32:
		return float32Value(database)
	case reflect.Float64:
		return float64Value(database)
	case reflect.Bool:
		return boolValue(database)
	case reflect.Int8:
		return intValue(database, 8, false)
	case reflect.Int16:
		return intValue(database, 16, false)
	case reflect.Int, reflect.Int32:
		return intValue(database, 32, false)
	case reflect.Int64:
		return intValue(database, 64, false)
	case reflect.Uint8:
		return intValue(database, 8, true)
	case reflect.Uint16:
		return intValue(database, 16, true)
	case reflect.Uint, reflect.Uint32:
		return intValue(database, 32, true)
	case reflect.Uint64:
		return intValue(database, 64, true)
	case reflect.Struct:
		{
			switch kind {
			case reflect.TypeOf(time.Time{}):
				return "TIMESTAMP"
			default:
				return ""
			}
		}
	default:
		return ""
	}
}

func intValue(database dbsys, bitSize uint8, isUnsigned bool) string {
	if database == Postgres {
		return ""
	}

	var val string

	switch bitSize {
	case 8:
		val = "SMALLINT"
	case 16:
		val = "MEDIUMINT"
	case 32:
		val = "INT"
	case 64:
		val = "BIGINT"
	default:
		return ""
	}

	if isUnsigned {
		val += " UNSIGNED"
	}
	return val
}

func isUnsigned(kind reflect.Kind) bool {
	switch kind {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
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
func (dbhelper *DBhelper) CreateTable(data interface{}, options ...*CreateOption) error {
	var option *CreateOption
	if len(options) > 0 {
		option = options[0]
	}
	return dbhelper.handleErrHook(dbhelper.create(data, option), "creating table ")
}

//Insert creates a table for struct
//Leave name empty to use the name of the struct
func (dbhelper *DBhelper) Insert(data interface{}, options ...*InsertOption) (*sql.Result, error) {
	var option *InsertOption
	if len(options) > 0 {
		option = options[0]
	}

	res, err := dbhelper.insert(data, option)
	return res, dbhelper.handleErrHook(err, "inserting table")
}

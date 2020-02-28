package godbhelper

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func strHasEmpty(arg ...string) bool {
	for _, s := range arg {
		if len(strings.Trim(s, " ")) == 0 {
			return true
		}
	}
	return false
}

func parseDSNstring(arg ...string) string {
	if len(arg) > 0 {
		dsn := "?" + strings.Join(arg, "&")
		return dsn
	}
	return ""
}

func parsePostgresString(arg ...string) string {
	if len(arg) > 0 {
		dsn := strings.Join(arg, " ")
		return dsn
	}
	return ""
}

func isPortValid(port uint16) bool {
	return port > 0 && port <= 65535
}

func stringArrToInterface(str []string) []interface{} {
	params := make([]interface{}, len(str))
	for i, p := range str {
		params[i] = p
	}
	return params
}

func strArrHas(arr []string, has string) bool {
	for _, a := range arr {
		if a == has {
			return true
		}
	}
	return false
}

func parsetTag(tagContent string) []string {
	if strings.Contains(tagContent, ",") {
		return strings.Split(tagContent, ",")
	}
	return []string{tagContent}
}

func strValueFromReflect(field reflect.Value) (string, string, error) {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), "0", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), "0", nil
	case reflect.String:
		return field.String(), "", nil
	case reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', 5, 64), "0.0", nil
	case reflect.Float32:
		return strconv.FormatFloat(field.Float(), 'f', 5, 32), "0.0", nil
	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), "0", nil
	case reflect.Struct:
		{
			switch field.Type() {
			case reflect.TypeOf(time.Time{}):
				t := field.Interface().(time.Time)
				if t.IsZero() {
					return "", "0", nil
				}
				return fmt.Sprintf("FROM_UNIXTIME(%d)", t.Unix()), "0", nil
			default:
				return "", "", errors.New("Struct " + field.Type().String() + " not supported")
			}
		}
	default:
		log.Printf("Kind %s not found!\n", field.Kind().String())
		return "", "", errors.New("Kind not supported")
	}
}

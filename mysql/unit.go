package mysql

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MapSS map[string]string

func StrVal(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

// ConvertValueToString
// @Description: 根据field的类型 转换成不同的value格式
// @param fieldType
// @param val
// @return string
func ConvertValueToString(fieldType string, val []byte) string {
	if val == nil {
		return "NULL"
	}
	if NumTypeContains(fieldType) {
		return string(val)
	} else if CharTypeContains(fieldType) {
		return "'" + escape(val) + "'"
	} else {
		return "'" + string(val) + "'"
	}
}
func CharTypeContains(str string) bool {
	intTypes := [...]string{"char", "text", "blob"}

	for i := range intTypes {
		if strings.Contains(str, intTypes[i]) {
			return true
		}
	}
	return false
}

func NumTypeContains(str string) bool {
	intTypes := [...]string{"int", "float", "double"}

	for i := range intTypes {
		if strings.Contains(str, intTypes[i]) {
			return true
		}
	}
	return false
}

func escape(source []byte) string {
	var j int
	if len(source) == 0 {
		return ""
	}
	tempStr := source
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
		case '\n':
			flag = true
			escape = '\n'
		case '\\':
			flag = true
			escape = '\\'
		case '\'':
			flag = true
			escape = '\''
		case '"':
			flag = true
			escape = '"'
		case '\032':
			flag = true
			escape = 'Z'
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}

func GetInsertSql(v interface{}) (string, []interface{}) {
	elem := reflect.ValueOf(v)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()
	tablename := normalize(elemType.Name())
	numfields := elem.NumField()
	insertfields := make([]string, 0, numfields)
	insertvalues := make([]interface{}, 0, numfields)
	for i := 0; i < numfields; i++ {
		curfield := elem.Field(i)
		curstruct := elemType.Field(i)
		field, val := parseField(curfield, curstruct)
		insertfields = append(insertfields, "`"+field+"`")
		insertvalues = append(insertvalues, val)
	}
	quotes := strings.Repeat("?,", len(insertfields))
	quotes = quotes[0 : len(quotes)-1]
	ret := fmt.Sprintf("insert into %s(%s) values(%s)", tablename, strings.Join(insertfields, ","), quotes)
	return ret, insertvalues
}

func parseField(v reflect.Value, s reflect.StructField) (string, interface{}) {
	fieldname := normalize(s.Name)
	var fieldvalue interface{}
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			fieldvalue = "1"
		} else {
			fieldvalue = "0"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldvalue = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldvalue = v.Uint()
	case reflect.Float32, reflect.Float64:
		fieldvalue = v.Float()
	case reflect.String:
		fieldvalue = v.String()
	}
	return fieldname, fieldvalue
}

func normalize(str string) string {
	return strings.ToLower(str)
}

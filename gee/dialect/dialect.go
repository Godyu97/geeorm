package dialect

import (
	"reflect"
)

//1.将 Go 语言的类型映射为数据库中的类型。

var dialectsM = map[string]Dialect{}

type Dialect interface {
	//将 Go 语言的类型转换为该数据库的数据类型。
	DataTypeOf(typ reflect.Value) string
	//返回某个表是否存在的 SQL 语句，参数是表名(table)。
	TableExistSQL(tableName string) (string, []any)
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsM[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsM[name]
	return dialect, ok
}

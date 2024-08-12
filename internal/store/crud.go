package store

import (
	"fmt"
	"reflect"
	"strings"
)

func Select[T any](query string, args ...interface{}) ([]T, error) {
	res := make([]T, 0, 10)
	err := db.Select(&res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type InsertOps struct {
	TableName string
}

func WithTableName(name string) func(*InsertOps) {
	return func(io *InsertOps) {
		io.TableName = name
	}
}

// Insert insert value to table
//
// T must be pointer to struct
func Insert[T any](value *T, tableName string) error {
	if value == nil {
		return nil
	}
	sb := &strings.Builder{}
	sb.WriteString("insert into ")
	sb.WriteString(tableName)
	sb.WriteByte('(')
	structValue := reflect.ValueOf(value).Elem()
	structType := structValue.Type()
	numFields := structValue.NumField()
	args := make([]any, 0, 10)
	marks := make([]string, 0, 10)
	for i := range numFields {
		args = append(args, structValue.Field(i).Interface())
		marks = append(marks, "?")
		column, ok := structType.Field(i).Tag.Lookup("db")
		if !ok {
			return fmt.Errorf("tag(db) doesn't exist.")
		}
		if i == numFields-1 {
			sb.WriteString(column)
		} else {
			sb.WriteString(column)
			sb.WriteByte(',')
		}
	}
	sb.WriteByte(')')
	sb.WriteString("values(")
	sb.WriteString(strings.Join(marks, ","))
	sb.WriteByte(')')
	_, err := db.Exec(sb.String(), args...)
	return err
}

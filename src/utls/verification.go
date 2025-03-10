package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func ConfirmDataHasAllFields(data any, fields []string) error {
	v := reflect.ValueOf(data)
	missing := make([]string, 0)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for _, fieldName := range fields {
		f := v.FieldByName(fieldName)
		if !f.IsValid() {
			continue
		}
		if f.Kind() == reflect.Ptr && f.IsNil() {
			missing = append(missing, fieldName)
		}
	}
	if len(missing) > 0 {
		return errors.New(fmt.Sprint("missing fields: ", missing))
	}
	return nil
}

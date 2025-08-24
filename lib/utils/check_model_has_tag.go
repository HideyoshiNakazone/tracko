package utils

import (
	"reflect"
	"strings"
)


func CheckModelHasTag(model interface{}, fieldName string, tagKey string, tagValue string) bool {
	t := reflect.TypeOf(model)

	fieldNames := strings.Split(fieldName, ".")
	if len(fieldNames) == 0 {
		return false
	}
	structField, ok := FindFieldByNameTree(&t, &fieldNames)
	if !ok {
		return false
	}
	return structField.Tag.Get(tagKey) == tagValue
}


func FindFieldByNameTree(t *reflect.Type, fieldNames *[]string) (*reflect.StructField, bool) {
	if len(*fieldNames) == 0 {
		return nil, false
	}

	currentFieldName := (*fieldNames)[0]
	for i := 0; i < (*t).NumField(); i++ {
		field := (*t).Field(i)
		if field.Tag.Get("mapstructure") != currentFieldName {
			continue
		}
		if len(*fieldNames) == 1 {
			return &field, true
		}
		if field.Type.Kind() == reflect.Struct {
			subFieldNames := (*fieldNames)[1:]
			return FindFieldByNameTree(&field.Type, &subFieldNames)
		}
		return nil, false
	}
	return nil, false
}

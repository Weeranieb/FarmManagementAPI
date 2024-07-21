package dbutil

import "reflect"

func StructToMap(input interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	flattenStruct(val, output)
	return output
}

func flattenStruct(val reflect.Value, output map[string]interface{}) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// If field is unexported, ignore it
		if !field.CanInterface() {
			continue
		}

		// Handle nested structs by recursion
		if field.Kind() == reflect.Struct {
			flattenStruct(field, output)
		} else {
			output[fieldName] = field.Interface()
		}
	}
}

package mysql

import (
	"reflect"
	"strings"
)

// buildQueryAndParamsList builds the placeholders for the new params and adds them to the params list
func buildQueryAndParamsList(params *[]interface{}, newParams interface{}) (queryPlaceholders string) {
	switch reflect.TypeOf(newParams).Kind() {
	case reflect.Slice:
		newParamsSlice := reflect.ValueOf(newParams)
		paramLen := newParamsSlice.Len()

		if paramLen > 0 {
			queryPlaceholders = strings.Repeat(",?", paramLen)
			queryPlaceholders = queryPlaceholders[1:] // remove first comma

			for i := 0; i < paramLen; i++ {
				*params = append(*params, newParamsSlice.Index(i).Interface())
			}
		}
	default:
		queryPlaceholders = "?"
		*params = append(*params, newParams)
	}

	return queryPlaceholders
}

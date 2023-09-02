package commons

import (
	"os"
	"reflect"
	"strings"
)

func FileExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && !f.IsDir()
}

func SearchStruct(s interface{}, query string) bool {
	query = strings.ToLower(query)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			if SearchStruct(f.Interface(), query) {
				return true
			}
		} else {
			if strings.Contains(strings.ToLower(f.String()), query) {
				return true
			}
		}
	}
	return false
}

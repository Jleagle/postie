package json

import (
	"fmt"
	"reflect"
)

// List creates a html list from an interface{}
func List(data interface{}) string {
	return _list(data, "")
}

func _list(data interface{}, returnString string) string {
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {

		switch value := data.(type) {
		case string:
		case int32, int64:
			return ""
		case interface{}:
			return List(value)
		default:
			fmt.Println("unknown")
			return ""
		}
	}

	return ""
}

package helper

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
)

func EncodeCursor(value interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", value)))
}

func DecodeCursor(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func GetFieldValue(obj interface{}, field string) interface{} {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val.FieldByName(camelToPascal(field)).Interface()
}

func camelToPascal(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})

	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

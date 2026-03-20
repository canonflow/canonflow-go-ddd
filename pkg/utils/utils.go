package utils

import (
	"reflect"
	"strings"
)

func SliceContains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

func Reverse[T any](s *[]T) {
	arr := *s
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func GetFieldTagName(obj interface{}, fieldName string) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(fieldName); ok {
		tag := f.Tag.Get("form")
		if tag == "" {
			tag = f.Tag.Get("json")
		}
		if tag != "" {
			// handle cases like `json:"anchor_id,omitempty"`
			return strings.Split(tag, ",")[0]
		}
	}
	return strings.ToLower(fieldName) // fallback
}

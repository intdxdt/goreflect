package goreflect

import (
	"fmt"
	"reflect"
)

func GetType(v any) string {
	return reflect.TypeOf(v).String()
}

func GetJSONTaggedFields(s interface{}) ([]string, error) {
	var tags = make([]string, 0, 32)
	var v = reflect.ValueOf(s)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a pointer to a struct")
	}

	v = v.Elem() // deref struct
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup("json"); ok {
			if tag == "-" {
				continue
			}
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

func GetValues(obj any, tags []string) ([]any, error) {
	var vals = make([]any, 0, len(tags))
	var v = reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a pointer to a struct")
	}

	v = v.Elem() // dereference the pointer to get the struct
	t := v.Type()

	var dict = make(map[string]bool, len(tags))
	for _, tag := range tags {
		dict[tag] = true
	}

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup("json"); ok && dict[tag] {
			vals = append(vals, v.Field(i).Interface())
		}
	}
	return vals, nil
}

func GetFieldReferences(obj any, fields []string) ([]interface{}, error) {
	var refs = make([]any, 0, len(fields))
	var v = reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a pointer to a struct")
	}

	v = v.Elem() // deref struct
	t := v.Type()

	var dict = make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		if tag, ok := field.Tag.Lookup("json"); ok {
			dict[tag] = i
		}
	}

	for _, field := range fields {
		if idx, ok := dict[field]; ok {
			refs = append(refs, v.Field(idx).Addr().Interface())
		} else {
			return nil, fmt.Errorf("field '%s' not found", field)
		}
	}

	return refs, nil
}

func FilterFieldReferences(fields []string, fieldRefMap map[string]any) ([]string, []any, error) {
	var cols = make([]string, 0, len(fields))
	var refs = make([]any, 0, len(fields))
	for _, field := range fields {
		if ref, ok := fieldRefMap[field]; ok {
			cols = append(cols, field)
			refs = append(refs, ref)
		} else {
			return nil, nil, fmt.Errorf("field '%s' not found", field)
		}
	}
	return cols, refs, nil
}

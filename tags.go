package resourcenames

import (
	"fmt"
	"reflect"
)

// structs that we do not want to descend into
// when walking structs
var knownStructs = [][2]string{
	{"time", "Time"},
}

func walkStruct(val reflect.Value, visited map[uintptr]bool, fn func(field reflect.StructField, value reflect.Value)) {
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return
		}

		ptr := val.Pointer()
		if visited[ptr] {
			return
		}
		visited[ptr] = true

		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		fn(field, fieldValue)

		// do not recurse into known structs
		for _, kind := range knownStructs {
			if fieldValue.Type().PkgPath() == kind[0] && fieldValue.Type().Name() == kind[1] {
				continue
			}
		}

		// recurse into nested struct
		switch fieldValue.Kind() {
		case reflect.Struct:
			walkStruct(fieldValue, visited, fn)
		case reflect.Pointer:
			if fieldValue.Elem().IsValid() && fieldValue.Elem().Kind() == reflect.Struct {
				walkStruct(fieldValue, visited, fn)
			}
		}
	}
}

func findFieldByTag(val reflect.Value, tag string) (reflect.Value, bool) {
	var result reflect.Value
	found := false

	walkStruct(val, map[uintptr]bool{}, func(field reflect.StructField, fieldValue reflect.Value) {
		if found {
			return
		}
		if field.Tag.Get("rns") == tag {
			result = fieldValue
			found = true
		}
	})

	return result, found
}

func (p NamePattern) Unmarshal(resourceName string, resource any) error {
	params, err := p.Parse(resourceName)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(resource)

	if val.Kind() != reflect.Pointer || val.IsNil() {
		return ErrNotAStruct
	}

	val = val.Elem()

	if val.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	walkStruct(val, map[uintptr]bool{}, func(field reflect.StructField, fieldValue reflect.Value) {
		rns := field.Tag.Get("rns")
		if rns == "" {
			return
		}

		segment, ok := params[rns]
		if !ok {
			return
		}

		if !fieldValue.CanSet() {
			return
		}

		segmentValue := reflect.ValueOf(segment)

		if segmentValue.Type().AssignableTo(fieldValue.Type()) {
			fieldValue.Set(segmentValue)
		} else if segmentValue.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(segmentValue.Convert(fieldValue.Type()))
		}
	})

	return nil
}

func (p NamePattern) Marshal(resource any) (string, error) {
	resourceName := ""

	val := reflect.ValueOf(resource)

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return "", ErrNotAStruct
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", ErrNotAStruct
	}

	for _, seg := range p.segments {

		if resourceName != "" {
			resourceName += "/"
		}

		if !seg.isParam {
			resourceName += seg.value
			continue
		}

		fieldValue, found := findFieldByTag(val, seg.value)
		if !found {
			return "", fmt.Errorf("%w: %q", ErrMissingSegment, seg.value)
		}

		var str string

		switch fieldValue.Kind() {
		case reflect.String:
			str = fieldValue.String()
		default:
			str = fmt.Sprint(fieldValue.Interface())
		}

		resourceName += str
	}

	return resourceName, nil
}

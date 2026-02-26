package resourcenames

import (
	"fmt"
	"reflect"
)

// MatchInto parses a given resource name with the name pattern.
// All discovered variable segment values will be set on the corresponding
// resource field.
// It chooses the field to set the value to by comparing the segment name
// with the matching value in the `rns:"<segment_name>"` tag.
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

	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		if rns := field.Tag.Get("rns"); rns != "" {
			if segment, ok := params[rns]; ok {
				if !fieldValue.CanSet() {
					continue
				}

				segmentValue := reflect.ValueOf(segment)
				if segmentValue.Type().AssignableTo(fieldValue.Type()) {
					fieldValue.Set(segmentValue)
				} else if segmentValue.Type().ConvertibleTo(fieldValue.Type()) {
					fieldValue.Set(segmentValue.Convert(fieldValue.Type()))
				}
			}
		}
	}

	return nil
}

func (p NamePattern) Marshal(resource any) (string, error) {
	resourceName := ""
	val := reflect.ValueOf(resource)

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return resourceName, ErrNotAStruct
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return resourceName, ErrNotAStruct
	}

	typ := val.Type()

	for _, segment := range p.segments {
		if !segment.isParam {
			resourceName += "/" + segment.value
			continue
		}

		found := false

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := val.Field(i)

			if field.Tag.Get("rns") != segment.value {
				continue
			}

			// Convert field value to string
			var str string

			switch fieldValue.Kind() {
			case reflect.String:
				str = fieldValue.String()
			default:
				str = fmt.Sprint(fieldValue.Interface())
			}

			resourceName += "/" + str
			found = true
			break
		}

		if !found {
			return "", fmt.Errorf("%w: %q", ErrMissingSegment, segment.value)
		}
	}
	return resourceName, nil
}

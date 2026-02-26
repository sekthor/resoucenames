package resourcenames

import "reflect"

func (p NamePattern) MatchInto(resourceName string, resource any) error {
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

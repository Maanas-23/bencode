package bencode

import (
	"fmt"
	"reflect"
)

// unmarshal populates the reflect.Value v with the data from rawData.
// v must be a settable value (a pointer or a settable field).
func unmarshal(rawData any, v reflect.Value) error {
	// If v is a pointer, set the value it points to.
	if v.Kind() == reflect.Pointer {
		// If the pointer is nil, create a new value for it to point to.
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		// Dereference the pointer.
		v = v.Elem()
	}

	// If rawData is nil, we can't do anything further.
	if rawData == nil {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		s, ok := rawData.(string)
		if !ok {
			return fmt.Errorf("bencode: cannot unmarshal %T into Go value of type string", rawData)
		}
		v.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, ok := rawData.(int64)
		if !ok {
			return fmt.Errorf("bencode: cannot unmarshal %T into Go value of type int64", rawData)
		}
		if v.OverflowInt(i) {
			return fmt.Errorf("bencode: value %d overflows Go value of type %s", i, v.Type())
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, ok := rawData.(int64)
		if !ok {
			return fmt.Errorf("bencode: cannot unmarshal %T into Go value of type uint64", rawData)
		}
		if i < 0 {
			return fmt.Errorf("bencode: cannot unmarshal negative value %d into unsigned Go type %s", i, v.Type())
		}
		if v.OverflowUint(uint64(i)) {
			return fmt.Errorf("bencode: value %d overflows Go value of type %s", i, v.Type())
		}
		v.SetUint(uint64(i))

	case reflect.Slice:
		rawSlice, ok := rawData.([]any)
		if !ok {
			return fmt.Errorf("bencode: cannot unmarshal %T into Go value of type slice", rawData)
		}
		slice := reflect.MakeSlice(v.Type(), len(rawSlice), len(rawSlice))
		for i, item := range rawSlice {
			if err := unmarshal(item, slice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(slice)

	case reflect.Struct:
		rawMap, ok := rawData.(map[string]any)
		if !ok {
			return fmt.Errorf("bencode: cannot unmarshal %T into Go value of type struct", rawData)
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			// Skip unexported fields.
			if field.PkgPath != "" {
				continue
			}

			tag := field.Tag.Get("bencode")
			if tag == "" {
				tag = field.Name // Default to field name if no tag
			}

			if rawValue, ok := rawMap[tag]; ok {
				if err := unmarshal(rawValue, v.Field(i)); err != nil {
					return err
				}
			}
		}

	case reflect.Map:
		rawMap, ok := rawData.(map[string]any)
		if !ok {
			return fmt.Errorf("bencode: cannot unmarshal %T into Go value of type map", rawData)
		}
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
		for key, rawValue := range rawMap {
			mapValue := reflect.New(v.Type().Elem()).Elem()
			if err := unmarshal(rawValue, mapValue); err != nil {
				return err
			}
			v.SetMapIndex(reflect.ValueOf(key), mapValue)
		}

	case reflect.Interface:
		if !v.IsNil() {
			currentType := v.Elem().Type()
			newValue := reflect.ValueOf(rawData)
			if !newValue.Type().AssignableTo(currentType) {
				return fmt.Errorf("bencode: cannot unmarshal %T into value of type %s", rawData, currentType)
			}
		}
		v.Set(reflect.ValueOf(rawData))

	default:
		return fmt.Errorf("bencode: unsupported type for unmarshaling: %s", v.Kind())
	}

	return nil
}

package goutils

import (
	"errors"
	"reflect"
)

// Copy from one to another
// struct ==> struct
// slice ==> slice
// slice ==> struct
func Copy(toValue interface{}, fromValue interface{}) (err error) {

	var (
		isSlice bool
		//isMap   bool
		sliceSize = 1
		to        = reflect.Indirect(reflect.ValueOf(toValue))
		from      = reflect.Indirect(reflect.ValueOf(fromValue))
	)

	// from value must not a zero Value.
	if !from.IsValid() {
		return errors.New("from value is a zero value")
	}

	// to value must addressable
	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	if from.Type().AssignableTo(to.Type()) {
		// if from's value is assignable to to's type, we can use set
		to.Set(from)
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			sliceSize = from.Len()
		}
	}

	for i := 0; i < sliceSize; i++ {
		var dest, source reflect.Value

		if isSlice {
			// dest
			dest = reflect.Indirect(reflect.New(toType))

			// source
			if from.Kind() == reflect.Slice {
				source = reflect.Indirect(from.Index(i))
			} else {
				source = reflect.Indirect(from)
			}
		} else {

			// source
			source = reflect.Indirect(from)

			// dest
			dest = reflect.Indirect(to)
		}

		for _, field := range deepFields(fromType) {
			name := field.Name

			if fromField := source.FieldByName(name); fromField.IsValid() {
				// try copy method
				var toMethod reflect.Value
				if dest.CanAddr() {
					toMethod = dest.Addr().MethodByName(name + "Copy")
				} else {
					toMethod = dest.MethodByName(name + "Copy")
				}
				if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
					toMethod.Call([]reflect.Value{fromField})
					// if has copy method, then don't try field name again
					continue
				}

				// try dest field
				if toField := dest.FieldByName(name); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField) {
							if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
								return err
							}
						}
					}

				} else {
					// try to set to method
					if dest.CanAddr() {
						toMethod = dest.Addr().MethodByName(name)
					} else {
						toMethod = dest.MethodByName(name)
					}

					if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
						toMethod.Call([]reflect.Value{fromField})
					}
				}

			}

		}

		// copy from method to field
		for _, field := range deepFields(toType) {
			name := field.Name

			var fromMethod reflect.Value
			if source.CanAddr() {
				fromMethod = source.Addr().MethodByName(name)
			} else {
				fromMethod = source.MethodByName(name)
			}

			if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
				if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						toField.Set(values[0])
					}
				}
			}

		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}

		}

	}

	return
}

// if the element type is Ptr, Slice or Map, then returns a type's element type.
func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// get all field, include embedded field
func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func set(to, from reflect.Value) bool {
	if !from.IsValid() || !to.CanSet() {
		return false
	}

	if to.Kind() == reflect.Ptr {
		if from.Kind() == reflect.Ptr && from.IsNil() {
			to.Set(reflect.Zero(to.Type()))
		} else if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}
		//to = to.Elem()
	}

	if from.Type().AssignableTo(to.Type()) {
		to.Set(from)
	} else if from.Type().ConvertibleTo(to.Type()) {
		to.Set(from.Convert(to.Type()))
	} else if from.Kind() == reflect.Ptr {
		return set(to, from.Elem())
	} else {
		return false
	}
	return true
}

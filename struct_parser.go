package marshal

import (
	"fmt"
	"reflect"
)

type StructParser struct {
	tagName       string
	valueResolver func(key string, val reflect.Value, fatherStructFrom map[string]interface{}, fathersVal reflect.Value) (reflect.Value, error)
}

func (parser *StructParser) Parse(from reflect.Value, into reflect.Value) error {
	concreteFrom, ok := from.Interface().(map[string]interface{})
	if !ok {
		return fmt.Errorf("struct can be parsed only from map")
	}

	for key, val := range concreteFrom {
		destinationField, err := parser.findFieldByTag(parser.tagName, key, into)
		if err != nil {
			return err
		}

		err = parser.storeInField(destinationField, reflect.ValueOf(val), key, concreteFrom, into)
		if err != nil {
			return err
		}
	}

	return nil
}

// storeInField storing the given val into the field.
// field 		- the field to store the value in.
// val 			- the value to store in the field.
// fieldKey 	- the key of the field in the map[string]interface{}
// fatherData 	- the entire map[string]interface{}
// fatherVal 	- the reflect.Value of the father object
func (parser *StructParser) storeInField(field reflect.Value, val reflect.Value, fieldKey string, fatherData map[string]interface{}, fatherVal reflect.Value) error {
	switch field.Kind() {
	case reflect.String, reflect.Int, reflect.Int32, reflect.Int64:
		return PrimitivesParser.Parse(val, field)

	case reflect.Struct:
		p := StructParser{tagName: parser.tagName, valueResolver: parser.valueResolver}
		return p.Parse(val, field)

	case reflect.Pointer:
		return parser.handlePointer(field, val, fieldKey, fatherData, fatherVal)

	case reflect.Interface:
		return parser.handleInterface(field, val, fieldKey, fatherData, fatherVal)
	}

	return fmt.Errorf("unknown kind %s", field.Kind().String())
}

func (parser *StructParser) findFieldByTag(tagName string, expectedTagValue string, typ reflect.Value) (reflect.Value, error) {
	// In case the given value is a pointer
	indirectType := reflect.Indirect(typ)

	for i := 0; i < indirectType.NumField(); i++ {
		if indirectType.Type().Field(i).Tag.Get(tagName) == expectedTagValue {
			return indirectType.Field(i), nil
		}
	}

	return reflect.Value{}, fmt.Errorf("unable to find field")
}

// handlePointer when the destination field is a pointer
// currently we support POINTERS TO STRUCTS only
//
// Checking if the pointer is to nil, if so - creating a new empty struct and pointing to it.
// Then letting the structParser to handle the assignment into the struct.
func (parser *StructParser) handlePointer(field reflect.Value, val reflect.Value, fieldKey string, fatherData map[string]interface{}, fatherVal reflect.Value) error {
	// Creating a new *T on the field value, so we'll have a place to set the value into. (in case it's nil)
	if field.IsZero() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	// Currently we assume that a pointer is ONLY to a struct
	p := StructParser{tagName: parser.tagName, valueResolver: parser.valueResolver}
	return p.Parse(val, field)
}

// handleInterface when field is an interface{}.
// We can't resolve the destination type we should store the value in cause it's an interface{}
// So we use 2 methods:
//
// Method1 - the father object can implement the UnionType interface. If so - we'll trigger's the fathers ParseValue() function to
// retrive the field's destination object.
//
// Method2 - if exists, trigger a function that was supplied in the parser's ctor, that can resolve the type for us.
//
// If both methods fails - we will ignore the interface{} and move it to the next field.
func (parser *StructParser) handleInterface(field reflect.Value, val reflect.Value, fieldKey string, fatherData map[string]interface{}, fatherVal reflect.Value) error {

	if parser.valueResolver != nil {
		// Creating a new field to put the value in
		newField, err := parser.valueResolver(fieldKey, val, fatherData, fatherVal)
		if err != nil {
			return err
		}

		// Storing the value on the new field
		err = parser.storeInField(newField, val, fieldKey, fatherData, fatherVal)
		if err != nil {
			return err
		}

		// Setting the value of the new field, on the current struct's field
		field.Set(newField)
		return nil
	}

	// PlanB - checking if the father have implemented the UnionType interface.

	return nil
}

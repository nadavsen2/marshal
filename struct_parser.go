package marshal

import (
	"fmt"
	"reflect"
)

type StructParser struct {
	config *Config
}

func NewStructParser(config *Config) *StructParser {
	return &StructParser{
		config: config,
	}
}

func (parser *StructParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	concreteFrom, ok := from.Interface().(map[string]interface{})
	if !ok {
		return fmt.Errorf("struct can be parsed only from map")
	}

	for key, val := range concreteFrom {
		destinationField, err := parser.findFieldByTag(parser.config.TagName, key, into)
		if err != nil {
			return err
		}

		newCtx := ParsingContext{
			KeyInFather: key,
			FromFather:  concreteFrom,
			FatherVal:   into,
		}
		err = parser.storeInField(reflect.ValueOf(val), destinationField, newCtx)
		if err != nil {
			return err
		}
	}

	return nil
}

// storeInField storing the given val into the field.
// from 			- the field to store the value in.
// into 			- the value to store in the field.
// fieldKey 		- the key of the field in the map[string]interface{}
// fatherData 		- the entire map[string]interface{}
// fatherVal 		- the reflect.Value of the father object
func (parser *StructParser) storeInField(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	switch from.Kind() {
	case reflect.String, reflect.Int, reflect.Int32, reflect.Int64:
		return NewResolverParser(parser.config).Parse(from, into, ctx)

	case reflect.Struct:
		return NewResolverParser(parser.config).Parse(from, into, ctx)

	case reflect.Map:
		return NewResolverParser(parser.config).Parse(from, into, ctx)

	case reflect.Pointer:
		return NewResolverParser(parser.config).Parse(from, into, ctx)

	case reflect.Interface:
		return NewResolverParser(parser.config).Parse(from, into, ctx)
		// return parser.handleInterface(from, into, ctx)
	}

	return fmt.Errorf("unknown kind %s", from.Kind().String())
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
func (parser *StructParser) handleInterface(field reflect.Value, val reflect.Value, ctx ParsingContext) error {
	if parser.config.ValueResolver != nil {
		// Creating a new field to put the value in
		newField, err := parser.config.ValueResolver(ctx)
		if err != nil {
			return err
		}

		// Storing the value on the new field
		err = parser.storeInField(newField, val, ctx)
		if err != nil {
			return err
		}

		// Setting the value of the new field, on the current struct's field
		field.Set(newField)
		return nil
	}

	// TODO:
	// PlanB - checking if the father have implemented the UnionType interface.

	return nil
}

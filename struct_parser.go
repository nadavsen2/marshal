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

		err = NewResolverParser(parser.config).Parse(reflect.ValueOf(val), destinationField, newCtx)
		if err != nil {
			return err
		}
	}

	return nil
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

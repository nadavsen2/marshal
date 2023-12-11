package marshal

import (
	"reflect"
)

type InterfaceParser struct {
	config *Config
}

func NewInterfaceParser(config *Config) *InterfaceParser {
	return &InterfaceParser{
		config: config,
	}
}

func (parser *InterfaceParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	var newField reflect.Value
	var err error

	if parser.config.ValueResolver != nil {
		// Creating a new field to put the value in
		newField, err = parser.config.ValueResolver(ctx)
		if err != nil {
			return err
		}
	}

	// PlanB - checking if the father have implemented the UnionType interface.
	// if so - we can ask him what's the current field actual data type.
	if ctx.FatherVal.Type().Implements(reflect.TypeOf(new(UnionTypeContainer)).Elem()) {
		container := ctx.FatherVal.Interface().(UnionTypeContainer)
		newField, err = container.ResolveType(ctx)
		if err != nil {
			return err
		}
	}

	// Storing the value on the new field
	err = NewResolverParser(parser.config).Parse(from, newField, ctx)
	if err != nil {
		return err
	}

	// Setting the value of the new field, on the current struct's field
	into.Set(newField)
	return nil
}

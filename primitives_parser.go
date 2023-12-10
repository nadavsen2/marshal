package marshal

import (
	"fmt"
	"reflect"
	"slices"
)

var primitives = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
	reflect.String,
}

func IsPrimitive(kind reflect.Kind) bool {
	return slices.Contains(primitives, kind)
}

var PrimitivesParser = &primitivesParser{}

type primitivesParser struct{}

func (parser *primitivesParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	if !into.CanSet() {
		return fmt.Errorf("cannot set primitive value")
	}

	switch t := from.Interface().(type) {
	case string:
		into.SetString(t)
	default:
		into.Set(from)
	}

	return nil
}

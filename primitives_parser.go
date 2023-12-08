package marshal

import (
	"fmt"
	"reflect"
)

var PrimitivesParser = &primitivesParser{}

type primitivesParser struct{}

func (parser *primitivesParser) Parse(from reflect.Value, into reflect.Value) error {
	if !into.CanSet() {
		return fmt.Errorf("cannot set primitive value")
	}

	into.Set(from)
	return nil
}

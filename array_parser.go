package marshal

import (
	"fmt"
	"reflect"
)

type ArrayParser struct {
	config *Config
}

func NewArrayParser(config *Config) *ArrayParser {
	return &ArrayParser{
		config: config,
	}
}

func (parser *ArrayParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	if from.Kind() != reflect.Array && from.Kind() != reflect.Slice {
		return fmt.Errorf("filling array from array only")
	}

	destinaionArr := into

	// If the destination array is nil, we create a new one
	if destinaionArr.IsZero() {
		destinaionArr = reflect.New(reflect.ArrayOf(from.Len(), into.Type().Elem())).Elem()
	}

	for i := 0; i < from.Len(); i++ {
		fromItem := from.Index(i)
		destinationItem := destinaionArr.Index(i)

		if err := NewResolverParser(parser.config).Parse(fromItem, destinationItem, ctx); err != nil {
			return fmt.Errorf("parsing array: %w", err)
		}
	}

	// Converting the result from Array (fixed size) into a slice.
	// e.g. [2]string to []string
	destinationSlice := destinaionArr.Slice(0, destinaionArr.Len())
	into.Set(destinationSlice)

	return nil
}

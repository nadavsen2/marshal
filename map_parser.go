package marshal

import (
	"fmt"
	"reflect"
)

type MapParser struct {
	config *Config
}

func NewMapParser(config *Config) *MapParser {
	return &MapParser{
		config: config,
	}
}

func (parser *MapParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	if from.Kind() != reflect.Map {
		return fmt.Errorf("filling a map only from map")
	}

	if from.Type().Key().Kind() != reflect.String || into.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("supporting only string key maps")
	}

	destinaionMap := into
	// If the destination map is nil, we create a new one
	if destinaionMap.IsZero() {
		mapType := reflect.MapOf(into.Type().Key(), into.Type().Elem())
		destinaionMap = reflect.MakeMap(mapType)
	}

	for _, keyVal := range from.MapKeys() {
		fromVal := from.MapIndex(keyVal)

		// The container for the value
		var intoVal = reflect.Indirect(reflect.New(into.Type().Elem()))

		// Letting the other parsers to deal with parsing T
		if err := NewResolverParser(parser.config).Parse(fromVal, intoVal, ctx); err != nil {
			return fmt.Errorf("parsing map: %w", err)
		}

		// Storing the parsed T inside the map
		destinaionMap.SetMapIndex(keyVal, intoVal)
	}

	into.Set(destinaionMap)

	return nil
}

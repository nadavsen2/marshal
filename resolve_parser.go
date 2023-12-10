package marshal

import (
	"fmt"
	"reflect"
)

type ResolverParser struct {
	config *Config
}

func NewResolverParser(config *Config) *ResolverParser {
	return &ResolverParser{
		config: config,
	}
}

func (parser *ResolverParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	switch into.Kind() {
	case reflect.String, reflect.Int, reflect.Int32, reflect.Int64:
		return PrimitivesParser.Parse(from, into, ctx)

	case reflect.Struct:
		return NewStructParser(parser.config).Parse(from, into, ctx)

	case reflect.Map:
		return NewMapParser(parser.config).Parse(from, into, ctx)

	case reflect.Pointer:
		return NewPointerParser(parser.config).Parse(from, into, ctx)

	case reflect.Interface:
		return NewInterfaceParser(parser.config).Parse(from, into, ctx)
	// 	return parser.handleInterface(field, val, fieldKey, fatherData, fatherVal)

	default:
		return fmt.Errorf("unsupported")
	}
}

package marshal

import "reflect"

type PointerParser struct {
	config *Config
}

func NewPointerParser(config *Config) *PointerParser {
	return &PointerParser{
		config: config,
	}
}

// handlePointer when the destination field is a pointer
// currently we support POINTERS TO STRUCTS only
//
// Checking if the pointer is to nil, if so - creating a new empty struct and pointing to it.
// Then letting the structParser to handle the assignment into the struct.
func (parser *PointerParser) Parse(from reflect.Value, into reflect.Value, ctx ParsingContext) error {
	// Creating a new *T on the field value, so we'll have a place to set the value into. (in case it's nil)
	if into.IsZero() {
		into.Set(reflect.New(into.Type().Elem()))
	}

	// Currently we assume that a pointer is ONLY to a struct
	return NewStructParser(parser.config).Parse(from, into, ctx)
}

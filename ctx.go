package marshal

import "reflect"

var EmptyContext = ParsingContext{}

type ParsingContext struct {

	// Represents the key within the father that we're trying to parse right now.
	KeyInFather string

	// The entire "from" data in which the father was constructed.
	FromFather map[string]interface{}

	// The reflect.Value of the father.
	FatherVal reflect.Value
}

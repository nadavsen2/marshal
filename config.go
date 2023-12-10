package marshal

import "reflect"

type Config struct {
	TagName       string
	ValueResolver func(ctx ParsingContext) (reflect.Value, error)
}

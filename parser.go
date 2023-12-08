package marshal

import "reflect"

type Parser interface {
	Parse(from reflect.Value, into reflect.Value) error
}

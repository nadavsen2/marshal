package marshal

import "reflect"

type UnionTypeContainer interface {
	ResolveType(ctx ParsingContext) (reflect.Value, error)
}

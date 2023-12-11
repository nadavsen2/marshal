package marshal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnionTypeWithContainerResolverString(t *testing.T) {
	expected := Container{
		Type: "string",
		Data: "data",
	}

	input := map[string]interface{}{
		"type": "string",
		"data": "data",
	}

	// Act:
	actual := Container{}

	structParser := NewStructParser(&Config{TagName: "marshal"})
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnionTypeWithContainerResolverStruct(t *testing.T) {
	expected := Container{
		Type: "struct",
		Data: &SomeData{
			D: 5,
		},
	}

	input := map[string]interface{}{
		"type": "struct",
		"data": map[string]interface{}{
			"d": 5,
		},
	}

	// Act:
	actual := Container{}

	structParser := NewStructParser(&Config{TagName: "marshal"})
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

type Container struct {
	Type string `marshal:"type"`
	Data any    `marshal:"data"`
}

type SomeData struct {
	D int `marshal:"d"`
}

func (st *Container) ResolveType(ctx ParsingContext) (reflect.Value, error) {
	switch ctx.FromFather["type"].(string) {
	case "string":
		return reflect.New(reflect.TypeOf("")).Elem(), nil
	case "struct":
		return reflect.ValueOf(&SomeData{}), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type")
	}
}

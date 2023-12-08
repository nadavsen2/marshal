package marshal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RootStruct struct {
	Name string `marshal:"name"`
	Type string `marshal:"type"`
	Strc DataB  `marshal:"struct"`
	Ptr  *DataA `marshal:"ptr"`
}

type DataA struct {
	A int `marshal:"a"`
}

type DataB struct {
	B          string `marshal:"b"`
	Underlying DataA  `marshal:"underlying"`
	Data       any    `marshal:"data"`
}

func TestStructWithUnionType(t *testing.T) {

	// Arrange:
	structParser := StructParser{
		tagName: "marshal",
		valueResolver: func(key string, val reflect.Value, fatherStructFrom map[string]interface{}, fatherVal reflect.Value) (reflect.Value, error) {
			switch fatherVal.Type() {
			case reflect.TypeOf(DataB{}):
				return reflect.ValueOf(&DataA{}), nil
			default:
				return reflect.Value{}, fmt.Errorf("unexpected value")
			}
		},
	}

	expected := RootStruct{
		Name: "test",
		Type: "A",
		Strc: DataB{
			B: "btest",
			Underlying: DataA{
				A: 555,
			},
			Data: &DataA{
				A: 1,
			},
		},
		Ptr: &DataA{
			A: 55,
		},
	}

	input := map[string]interface{}{
		"name": "test",
		"type": "A",
		"struct": map[string]interface{}{
			"b": "btest",
			"underlying": map[string]interface{}{
				"a": 555,
			},
			"data": map[string]interface{}{
				"a": 1,
			},
		},
		"ptr": map[string]interface{}{
			"a": 55,
		},
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual))

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

package marshal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RootStruct struct {
	Name                 string                 `marshal:"name"`
	Type                 string                 `marshal:"type"`
	Strc                 DataB                  `marshal:"struct"`
	Ptr                  *DataA                 `marshal:"ptr"`
	RegularMap           map[string]int         `marshal:"regular_map"`
	MapWithStruct        map[string]DataA       `marshal:"map_with_struct"`
	MapWithPointer       map[string]*DataA      `marshal:"map_with_pointer"`
	MapWithInterface     map[string]interface{} `marshal:"map_with_interface"`
	RegularArray         []string               `marshal:"regular_array"`
	ArrayWithStruct      []DataA                `marshal:"array_with_struct"`
	ArrayWithPointer     []*DataA               `marshal:"array_with_pointer"`
	UnionTypeToPrimitive any                    `marshal:"union_type_to_primitive"`
}

type DataA struct {
	A int `marshal:"a"`
}

type DataB struct {
	B          string `marshal:"b"`
	Underlying DataA  `marshal:"underlying"`
	Data       any    `marshal:"data"`
}

func TestUnionTypeToPrimitive(t *testing.T) {
	config := &Config{
		TagName: "marshal",
		ValueResolver: func(ctx ParsingContext) (reflect.Value, error) {
			return reflect.New(reflect.TypeOf("")).Elem(), nil
		},
	}
	structParser := NewStructParser(config)

	expected := RootStruct{
		UnionTypeToPrimitive: "test",
	}

	input := map[string]interface{}{
		"union_type_to_primitive": "test",
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestStructWithString(t *testing.T) {
	config := &Config{
		TagName: "marshal",
	}
	structParser := NewStructParser(config)

	expected := RootStruct{
		Name: "test",
		Type: "A",
	}

	input := map[string]interface{}{
		"name": "test",
		"type": "A",
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestRegularArray(t *testing.T) {
	config := &Config{
		TagName: "marshal",
	}

	structParser := NewStructParser(config)

	expected := RootStruct{
		RegularArray: []string{"a", "b"},
	}

	input := map[string]interface{}{
		"regular_array": []interface{}{"a", "b"},
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestArrayWithStruct(t *testing.T) {
	config := &Config{
		TagName: "marshal",
	}

	structParser := NewStructParser(config)
	expected := RootStruct{
		ArrayWithStruct: []DataA{
			{
				A: 1,
			},
			{
				A: 2,
			},
		},
	}

	input := map[string]interface{}{
		"array_with_struct": []interface{}{
			map[string]interface{}{
				"a": 1,
			},
			map[string]interface{}{
				"a": 2,
			},
		},
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestArrayWithPointer(t *testing.T) {
	config := &Config{
		TagName: "marshal",
	}

	structParser := NewStructParser(config)
	expected := RootStruct{
		ArrayWithPointer: []*DataA{
			{
				A: 1,
			},
			{
				A: 2,
			},
		},
	}

	input := map[string]interface{}{
		"array_with_pointer": []interface{}{
			map[string]interface{}{
				"a": 1,
			},
			map[string]interface{}{
				"a": 2,
			},
		},
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}

func TestMapWithStructs(t *testing.T) {
	// Arrange:

	input := map[string]interface{}{
		"map_with_struct": map[string]interface{}{
			"test1": map[string]interface{}{
				"a": 1,
			},
			"test2": map[string]interface{}{
				"a": 2,
			},
		},
	}

	expected := RootStruct{
		MapWithStruct: map[string]DataA{
			"test1": {
				A: 1,
			},
			"test2": {
				A: 2,
			},
		},
	}

	config := &Config{TagName: "marshal"}
	structParser := NewStructParser(config)

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMapWithPointers(t *testing.T) {
	// Arrange:

	input := map[string]interface{}{
		"map_with_pointer": map[string]interface{}{
			"test1": map[string]interface{}{
				"a": 1,
			},
			"test2": map[string]interface{}{
				"a": 2,
			},
		},
	}

	expected := RootStruct{
		MapWithPointer: map[string]*DataA{
			"test1": {
				A: 1,
			},
			"test2": {
				A: 2,
			},
		},
	}

	config := &Config{TagName: "marshal"}
	structParser := NewStructParser(config)

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMapWithInterface(t *testing.T) {
	// Arrange:
	input := map[string]interface{}{
		"map_with_interface": map[string]interface{}{
			"test1": map[string]interface{}{
				"a": 1,
			},
			"test2": map[string]interface{}{
				"a": 2,
			},
		},
	}

	expected := RootStruct{
		MapWithInterface: map[string]interface{}{
			"test1": &DataA{
				A: 1,
			},
			"test2": &DataA{
				A: 2,
			},
		},
	}

	config := &Config{
		TagName: "marshal",
		ValueResolver: func(ctx ParsingContext) (reflect.Value, error) {
			return reflect.ValueOf(&DataA{}), nil
		},
	}

	structParser := NewStructParser(config)

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestAllTogether(t *testing.T) {

	// Arrange:
	config := &Config{
		TagName: "marshal",
		ValueResolver: func(ctx ParsingContext) (reflect.Value, error) {
			switch ctx.FatherVal.Type() {
			case reflect.TypeOf(DataB{}):
				return reflect.ValueOf(&DataA{}), nil
			default:
				return reflect.Value{}, fmt.Errorf("unexpected value")
			}
		},
	}
	structParser := NewStructParser(config)

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
		RegularMap: map[string]int{
			"a": 1,
			"b": 2,
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
		"regular_map": map[string]int{
			"a": 1,
			"b": 2,
		},
	}

	// Act:
	actual := RootStruct{}
	err := structParser.Parse(reflect.ValueOf(input), reflect.ValueOf(&actual), EmptyContext)

	// Assert:
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

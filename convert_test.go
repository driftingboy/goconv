package goconv_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/driftingboy/goconv"
	"github.com/stretchr/testify/assert"
)

func Example_GetGenericValue() {
	// use default TypeDirect
	dc := goconv.NewDTOConverter()

	// get Value of this type
	// int
	i := new(int)
	err := dc.ConvertReflect("int", "10", i)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*i)
	fmt.Println(reflect.TypeOf(*i).String())
	//Output:
	//10
	//int
}

func Test_TypeDirect(t *testing.T) {
	// direct := goconv.NewTypeDirectWithConfig(false, map[goconv.Typ][]string{})
	// dc := goconv.NewDTOConverter(goconv.WithTypeDirect(direct))
	// use default TypeDirect
	dc := goconv.NewDTOConverter()
	// add some aliases
	goconv.DefaultTypeDirect.AddTypeAliases(goconv.Int, "interage")

	i := new(int)
	err := dc.ConvertReflect("interage", "100", i)
	if assert.NoError(t, err) {
		assert.Equal(t, 100, *i)
	}
}

func Test_GetGenericValue(t *testing.T) {

	dc := goconv.NewDTOConverter()

	tests := []struct {
		name string

		typ        string
		value      string
		wantResult interface{}
	}{
		{name: "int", typ: "int", value: "10", wantResult: interface{}(10)},
		{name: "int32", typ: "int32", value: "10", wantResult: interface{}(int32(10))},
		{name: "int64", typ: "int64", value: "10", wantResult: interface{}(int64(10))},
		{name: "float32", typ: "float32", value: "1.1", wantResult: interface{}(float32(1.1))},
		{name: "float64", typ: "float64", value: "1.1", wantResult: interface{}(float64(1.1))},
		{name: "string", typ: "string", value: "abc", wantResult: interface{}("abc")},
		{name: "bool", typ: "bool", value: "true", wantResult: interface{}(true)},
		{name: "date", typ: "date", value: "2022-01-01", wantResult: interface{}(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local))},
		{name: "datetime", typ: "datetime", value: "2022-01-01 12:00:00", wantResult: interface{}(time.Date(2022, time.January, 1, 12, 0, 0, 0, time.Local))},
		{name: "ints", typ: "ints", value: "1,2,3,4", wantResult: interface{}([]int{1, 2, 3, 4})},
		{name: "float32s", typ: "float32s", value: "1,2,3,4", wantResult: interface{}([]float32{1, 2, 3, 4})},
		{name: "float64s", typ: "float64s", value: "1,2,3,4", wantResult: interface{}([]float64{1, 2, 3, 4})},
		{name: "bools", typ: "bools", value: "true,false,true", wantResult: interface{}([]bool{true, false, true})},
		{name: "strings", typ: "strings", value: "1,2,3,4", wantResult: interface{}([]string{"1", "2", "3", "4"})},
		{
			name:  "dates",
			typ:   "dates",
			value: "2022-01-01,2022-01-02",
			wantResult: interface{}([]time.Time{
				time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local),
				time.Date(2022, time.January, 2, 0, 0, 0, 0, time.Local),
			}),
		},
		{
			name:  "dates",
			typ:   "datetimes",
			value: "2022-01-01 12:00:00,2022-01-02 12:00:00",
			wantResult: interface{}([]time.Time{
				time.Date(2022, time.January, 1, 12, 0, 0, 0, time.Local),
				time.Date(2022, time.January, 2, 12, 0, 0, 0, time.Local),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := dc.Convert(tt.typ, tt.value)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}

	t.Log(time.Date(2022, time.January, 1, 12, 30, 0, 0, time.Local).Format("2006-01-02 15:04:05"))
	t.Log(time.Date(2022, time.January, 1, 12, 30, 0, 0, time.Local).UTC())
	t.Log(time.Date(2022, time.January, 1, 12, 30, 0, 0, time.Local).Unix())
}

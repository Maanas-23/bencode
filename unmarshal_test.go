package bencode

import (
	"reflect"
	"testing"
)

type unmarshalTest struct {
	name    string
	in      string
	out     any
	want    any
	wantErr bool
}

func ptr[T any](v T) *T {
	return &v
}

var unmarshalTests = []unmarshalTest{
	{
		name: "Simple String",
		in:   "4:spam",
		out:  new(string),
		want: ptr("spam"),
	},
	{
		name: "Simple Integer",
		in:   "i42e",
		out:  new(int),
		want: ptr(42),
	},
	{
		name: "Negative Integer",
		in:   "i-42e",
		out:  new(int),
		want: ptr(-42),
	},
	{
		name: "Simple List",
		in:   "l4:spami42ee",
		out:  new([]any),
		want: &[]any{"spam", int64(42)},
	},
	{
		name: "Simple Dictionary",
		in:   "d3:foo3:bar5:helloi42ee",
		out:  new(map[string]any),
		want: &map[string]any{"foo": "bar", "hello": int64(42)},
	},
	{
		name: "Nested Struct",
		in:   "d4:dictd3:key5:valuee4:listli1ei2ei3eee",
		out: &struct {
			Dict struct {
				Key string `bencode:"key"`
			} `bencode:"dict"`
			List []int `bencode:"list"`
		}{},
		want: &struct {
			Dict struct {
				Key string `bencode:"key"`
			} `bencode:"dict"`
			List []int `bencode:"list"`
		}{
			Dict: struct {
				Key string `bencode:"key"`
			}{Key: "value"},
			List: []int{1, 2, 3},
		},
	},
	{
		name:    "Integer Overflow",
		in:      "i9223372036854775807e",
		out:     new(int8),
		wantErr: true,
	},
	{
		name:    "Unsigned Integer Overflow",
		in:      "i256e",
		out:     new(uint8),
		wantErr: true,
	},
	{
		name:    "Negative to Unsigned",
		in:      "i-1e",
		out:     new(uint),
		wantErr: true,
	},
	{
		name:    "Type Mismatch String to Int",
		in:      "4:spam",
		out:     new(int),
		wantErr: true,
	},
	{
		name:    "Type Mismatch Int to String",
		in:      "i42e",
		out:     new(string),
		wantErr: true,
	},
	{
		name:    "Type Mismatch List to Struct",
		in:      "li1ee",
		out:     &struct{}{},
		wantErr: true,
	},
	{
		name:    "Type Mismatch Dict to Slice",
		in:      "de",
		out:     &[]any{},
		wantErr: true,
	},
	{
		name: "Unmarshal into nil interface",
		in:   "i42e",
		out: func() any {
			var i any
			return &i
		}(),
		want: ptr(int64(42)),
	},
	{
		name: "Unmarshal into non-nil interface",
		in:   "i42e",
		out: func() any {
			var i any = new(string)
			return &i
		}(),
		wantErr: true,
	},
	{
		name: "Unmarshal into struct with unexported fields",
		in:   "d3:foo3:bar3:baz3:quxe",
		out: &struct {
			Foo string `bencode:"foo"`
			baz string `bencode:"baz"`
		}{},
		want: &struct {
			Foo string `bencode:"foo"`
			baz string `bencode:"baz"`
		}{
			Foo: "bar",
		},
	},
}

func TestUnmarshal(t *testing.T) {
	for _, tc := range unmarshalTests {
		t.Run(tc.name, func(t *testing.T) {
			err := Unmarshal([]byte(tc.in), tc.out)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				// Dereference the pointer to get the actual value
				val := reflect.ValueOf(tc.out).Elem().Interface()
				want := reflect.ValueOf(tc.want).Elem().Interface()
				if !reflect.DeepEqual(val, want) {
					t.Errorf("Unmarshal() got = %#v, want %#v", val, want)
				}
			}
		})
	}
}

func TestUnmarshalInvalidInput(t *testing.T) {
	var v any
	err := Unmarshal([]byte("i42e"), v)
	if err == nil {
		t.Error("expected an error for nil interface")
	}

	var s string
	err = Unmarshal([]byte("i42e"), s)
	if err == nil {
		t.Error("expected an error for non-pointer")
	}
}

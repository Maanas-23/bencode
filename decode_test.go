package bencode

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshalGeneric(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		want any
	}{
		{name: "Simple String", in: "4:spam", want: "spam"},
		{name: "Simple Integer", in: "i42e", want: int64(42)},
		{name: "Negative Integer", in: "i-42e", want: int64(-42)},
		{name: "Simple List", in: "l4:spami42ee", want: []any{"spam", int64(42)}},
		{name: "Simple Dictionary", in: "d3:foo3:bar5:helloi42ee", want: map[string]any{"foo": "bar", "hello": int64(42)}},
		{name: "Nested Structures",
			in: "d4:dictd3:key5:valuee4:listli1ei2ei3eee",
			want: map[string]any{
				"dict": map[string]any{"key": "value"},
				"list": []any{int64(1), int64(2), int64(3)},
			},
		},
		{name: "Empty string", in: "0:", want: ""},
		{name: "Empty list", in: "le", want: []any{}},
		{name: "Empty dictionary", in: "de", want: map[string]any{}},
		{name: "Binary Data in String", in: "2:\x01\x02", want: "\x01\x02"},
		{name: "Binary Data in Dict Key", in: "d2:\x01\x02i1ee", want: map[string]any{"\x01\x02": int64(1)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			err := Unmarshal([]byte(tc.in), &got)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Unmarshal() got = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestUnmarshalError(t *testing.T) {
	testCases := []struct {
		name string
		in   string
	}{
		{name: "Malformed String Length", in: "5:abc"},
		{name: "Malformed Integer No End", in: "i42"},
		{name: "Malformed List No End", in: "l4:spam"},
		{name: "Malformed Dictionary No End", in: "d3:foo3:bar"},
		{name: "Invalid Start Token", in: "x"},
		{name: "Lone End Token", in: "e"},
		{name: "Integer with non-digit chars", in: "i42a2e"},
		{name: "Dictionary with non-string key", in: "di1e3:fooee"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got any
			err := Unmarshal([]byte(tc.in), &got)
			if err == nil {
				t.Fatalf("Expected an error but got nil")
			}
		})
	}
}

func TestDecoderConsecutive(t *testing.T) {
	d := NewDecoder(strings.NewReader("i1ei2e4:spam"))
	var i int
	var s string

	err := d.Decode(&i)
	if err != nil || i != 1 {
		t.Fatalf("Expected to decode 1, got %d with err: %v", i, err)
	}

	err = d.Decode(&i)
	if err != nil || i != 2 {
		t.Fatalf("Expected to decode 2, got %d with err: %v", i, err)
	}

	err = d.Decode(&s)
	if err != nil || s != "spam" {
		t.Fatalf("Expected to decode 'spam', got %s with err: %v", s, err)
	}

	err = d.Decode(&i) // Should be EOF
	if err != io.EOF {
		t.Fatalf("Expected io.EOF, got %v", err)
	}
}

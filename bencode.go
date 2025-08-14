package bencode

import (
	"bytes"
	"io"
	"reflect"
)

// Unmarshal decodes the given Bencoded data into the given value.
func Unmarshal(data []byte, v any) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

// A Decoder reads and decodes Bencode values from an input stream.
type Decoder struct {
	r *reader
}

// NewDecoder returns a new decoder that reads from r.
//
// The decoder introduces its own buffering and may read data from r beyond the Bencode values requested.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: newReader(r)}
}

// Decode reads the next Bencode-encoded value from its
// input and returns it as an any
func (d *Decoder) Decode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}

	rawData, err := d.r.decode()
	if err != nil {
		return err
	}

	return unmarshal(rawData, rv)
}

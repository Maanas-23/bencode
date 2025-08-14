package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// reader is a buffered reader that provides methods for decoding bencode values.
type reader struct {
	r *bufio.Reader
}

// newReader creates a new reader from an io.Reader.
// If the reader is already a *bufio.Reader, it will be used directly.
func newReader(r io.Reader) *reader {
	if br, ok := r.(*bufio.Reader); ok {
		return &reader{r: br}
	}
	return &reader{r: bufio.NewReader(r)}
}

func (r *reader) decode() (any, error) {
	// Look at the first byte to determine the data type of value
	b, err := r.r.ReadByte()
	if err != nil {
		return nil, err
	}

	// Put the byte back so the respective parsing function can consume it.
	if err := r.r.UnreadByte(); err != nil {
		return nil, err
	}

	switch b {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return r.decodeString()
	case 'i':
		return r.decodeInt()
	case 'l':
		return r.decodeList()
	case 'd':
		return r.decodeDict()
	default:
		return nil, errors.New("bencode: invalid or unsupported type character")
	}
}

// decodeString parses a string from the reader.
// Format: <length>:<contents>
func (r *reader) decodeString() (string, error) {
	lengthStr, err := r.r.ReadString(':')
	if err != nil {
		if err == io.EOF {
			return "", errors.New("bencode: invalid string format, unexpected EOF")
		}
		return "", fmt.Errorf("bencode: invalid string format: %w", err)
	}
	lengthStr = lengthStr[:len(lengthStr)-1] // Remove the trailing ':'

	length, err := strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("bencode: invalid string length: %w", err)
	}

	contents := make([]byte, length)
	_, err = io.ReadFull(r.r, contents)
	if err != nil {
		return "", fmt.Errorf("bencode: failed to read string contents: %w", err)
	}

	return string(contents), nil
}

// decodeInt parses an integer from the reader.
// Format: i<integer>e
func (r *reader) decodeInt() (int64, error) {
	if b, err := r.r.ReadByte(); err != nil || b != 'i' {
		return 0, errors.New("bencode: expected 'i' at start of integer")
	}

	intStr, err := r.r.ReadString('e')
	if err != nil {
		return 0, fmt.Errorf("bencode: invalid integer format, could not find 'e': %w", err)
	}
	intStr = intStr[:len(intStr)-1] // Remove the trailing 'e'

	val, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("bencode: invalid integer value: %w", err)
	}

	return val, nil
}

// decodeList parses a list of Bencode values from the reader.
// Format: l<value1><value2>...e
func (r *reader) decodeList() ([]any, error) {
	if b, err := r.r.ReadByte(); err != nil || b != 'l' {
		return nil, errors.New("bencode: expected 'l' at start of list")
	}

	list := make([]any, 0)
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if err := r.r.UnreadByte(); err != nil {
			return nil, err
		}

		if b == 'e' {
			_, _ = r.r.ReadByte() // Consume the 'e'
			break
		}

		item, err := r.decode()
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

// decodeDict parses a dictionary of Bencode values from the reader.
// Format: d<key1><value1><key2><value2>...e
func (r *reader) decodeDict() (map[string]any, error) {
	if b, err := r.r.ReadByte(); err != nil || b != 'd' {
		return nil, errors.New("bencode: expected 'd' at start of dictionary")
	}

	dict := make(map[string]any)
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if err := r.r.UnreadByte(); err != nil {
			return nil, err
		}

		if b == 'e' {
			_, _ = r.r.ReadByte() // Consume the 'e'
			break
		}

		key, err := r.decodeString()
		if err != nil {
			return nil, fmt.Errorf("bencode: dictionary key must be a string: %w", err)
		}

		value, err := r.decode()
		if err != nil {
			return nil, err
		}
		dict[key] = value
	}

	return dict, nil
}

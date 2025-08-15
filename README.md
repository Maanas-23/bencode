## Bencode

*bencode* is a robust and easy-to-use Go package for decoding data in Bencode format, the encoding method used by the BitTorrent file-sharing protocol. This implementation is modeled after the standard library's `encoding/json` package, providing a familiar API.

Bencode supports four data types: strings, integers, lists, and dictionaries.

### Installation

```bash
go get github.com/maanas-23/bencode
```

### Usage

The `bencode` package makes it simple to unmarshal Bencoded byte slices into Go types.

#### Unmarshaling to Structs

For structured data, you can unmarshal Bencoded dictionaries directly into Go structs. Use the `bencode` struct tag to map dictionary keys to struct fields.

```go
package main

import (
	"fmt"
	"log"

	"github.com/maanas-23/bencode"
)

func main() {
	// A Bencoded dictionary: d3:foo3:bar5:counti42ee
	data := []byte("d3:foo3:bar5:counti42ee")

	type Info struct {
		Foo   string `bencode:"foo"`
		Count int    `bencode:"count"`
	}

	var info Info
	err := bencode.Unmarshal(data, &info)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Foo: %s\n", info.Foo)
	fmt.Printf("Count: %d\n", info.Count)
}
```

#### Unmarshaling to a Generic Interface

You can also unmarshal data into a generic `any` type (`interface{}`) if the structure is unknown or dynamic. Bencoded dictionaries are unmarshaled into `map[string]any`, and lists are unmarshaled into `[]any`.

```go
package main

import (
	"fmt"
	"log"

	"github.com/maanas-23/bencode"
)

func main() {
	// A Bencoded dictionary: d3:foo3:bar5:helloi42ee
	data := []byte("d3:foo3:bar5:helloi42ee")

	var v any
	err := bencode.Unmarshal(data, &v)
	if err != nil {
		log.Fatal(err)
	}

	// Type assert the result to a map
	dataMap, ok := v.(map[string]any)
	if !ok {
		log.Fatal("Expected a dictionary")
	}

	fmt.Printf("foo: %s\n", dataMap["foo"])
	fmt.Printf("hello: %d\n", dataMap["hello"])
}
```

## Bencode

*bencode* is a robust and easy-to-use Go package for encoding and decoding Bencode-formatted data, the encoding method used by the BitTorrent file-sharing protocol. This implementation is modeled after the standard library's encoding/json package.

Bencode supports four data types: strings, integers, lists, and dictionaries.

### Installation

```bash
go get github.com/Maanas-23/bencode
```

### Quick Start

You can easily parse a Bencoded byte slice using `bencode.Unmarshal`.

```go
package main

import (
	"fmt"
	"log"

	"github.com/Maanas-23/bencode"
)

func main() {
	// A Bencoded dictionary: d3:foo3:bar5:helloi42ee
	data := []byte("d3:foo3:bar5:helloi42ee")

	v, err := bencode.Unmarshal(data)
	if err != nil {
		log.Fatal(err)
	}

	// Type assert the result to a map
	dataMap, ok := v.(map[string]interface{})
	if !ok {
		log.Fatal("Expected a dictionary")
	}

	fmt.Printf("foo: %s\n", dataMap["foo"])
	fmt.Printf("hello: %d\n", dataMap["hello"])
}
```


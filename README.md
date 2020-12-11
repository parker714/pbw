# pbw
pbw is a HTTP web framework written in Go (Golang). 

[![GitHub license](https://img.shields.io/github/license/parker714/pbw)](https://github.com/parker714/pbw/blob/main/LICENSE)

## Supported Go versions

pbw is available as a [Go module](https://github.com/golang/go/wiki/Modules).

- 1.12+ 

## Feature Overview

- Route parameter binding
- Group APIs
- Extensible middleware framework
- Data binding for JSON and form payload

## Benchmarks

- 1

### Installation

```sh
go get github.com/parker714/pbw
```

### Example

```go
package main

import (
	"fmt"
	"github.com/parker714/pbw"
	"net/http"
)

func main() {
	e := pbw.New()
    
	e.Use(pbw.Recovery())
	e.Use(func(context pbw.Context) {
		fmt.Println("global middleware")
	})

	e.GET("/hello", func(c pbw.Context) {
		c.Data(http.StatusOK, []byte("hello"))
	})

	r1 := e.Group("/vip")
	{
		r1.Use(func(c pbw.Context) {
			fmt.Println("/vip middleware")
		})

		r1.GET("/user", func(c pbw.Context) {
			c.Data(http.StatusOK, []byte("user pb"))
		})
	}

	if err := e.Run(":8999"); err != nil {
		fmt.Printf("engine: run(:8999) fail, %s", err)
	}
}
```

## License

[Apache License 2.0](https://github.com/parker714/pbw/blob/main/LICENSE)
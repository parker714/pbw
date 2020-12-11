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

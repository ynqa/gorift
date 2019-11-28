package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"

	"github.com/gorift/gorift/pkg/resolve"
	"github.com/gorift/gorift/pkg/server"
)

func main() {
	resolver, err := resolve.NewDefaultResolver()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	report, err := resolver.Lookup(resolve.Request{
		Host: server.Host("localhost"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%# v", pretty.Formatter(report))
}

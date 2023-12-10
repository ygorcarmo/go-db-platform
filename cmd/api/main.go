package main

import (
	"custom-db-platform/internal/server"
	"fmt"
)

func main() {
	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}
}

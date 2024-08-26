//go:build dev
// +build dev

package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Unable to load .env")
	}
}

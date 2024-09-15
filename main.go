package main

import (
	"log"
	"os"

	"db-platform/api"

	"db-platform/storage"
)

func main() {

	store := storage.NewMySQLStorage(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_NAME"),
	)

	listenAddr := os.Getenv("LISTEN_ADDR")
	s := api.NewServer(listenAddr, store)
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ygorcarmo/db-platform/api"
	"github.com/ygorcarmo/db-platform/storage"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Unable to load .env")
	}

	store := storage.NewMySQLStorage(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_NAME"),
	)

	listenAddr := os.Getenv("LISTEN_ADDR")
	s := api.NewServer(listenAddr, store)

	log.Fatal(s.Start())

}

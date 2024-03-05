package main

import (
	"custom-db-platform/src/db"
	"custom-db-platform/src/server"
	"custom-db-platform/src/views"
	"fmt"
)

func init() {
	db.ConnectDB()
	views.LoadTemplates()
}

func main() {

	// hashParams := models.HashParams{
	// 	Memory:      64 * 1024,
	// 	Iterations:  3,
	// 	Parallelism: 2,
	// 	SaltLength:  16,
	// 	KeyLength:   32,
	// }

	// // sign up
	// encodedHash, err := hashParams.HashPasword("test")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(encodedHash)

	// // sign in
	// match, err := hashParams.ComparePasswordAndHash("test", encodedHash)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Match: %v\n", match)

	webServer := server.NewServer()
	err := webServer.ListenAndServe()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Cannot Start Server")
	}

}

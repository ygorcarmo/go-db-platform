package server

import (
	"fmt"
	"net/http"
	"time"
)

var port = 3000

type Server struct {
	port int
}

func NewServer() *http.Server {

	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Printf("Server is running on http://localhost:%d\n", port)

	return server
}

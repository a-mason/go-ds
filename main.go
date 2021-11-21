package main

import (
	"log"
	"github.com/a-mason/go-ds/internal/server"
)

func main() {
	const port = ":8080"
	log.Println("Starting server on port: " + port)
	srv := server.NewHttpServer(port)
	log.Fatal(srv.ListenAndServe())
}

package main

import (
	"github.com/a-mason/go-ds/internal/server"
	"log"
)

func main() {
	const port = ":8080"
	log.Println("Starting server on port: " + port)
	srv := server.NewHttpServer(port)
	log.Fatal(srv.ListenAndServe())
}

package main

import (
	"github.com/edwinfinlay/goland-api/internal/db"
	"github.com/edwinfinlay/goland-api/internal/todo"
	"github.com/edwinfinlay/goland-api/internal/transport"
	"log"
)

func main() {
	d, err := db.New("postgres", "example", "postgres", "localhost", 5432)
	if err != nil {
		log.Fatal(err)
	}

	svc := todo.NewService(d)
	server := transport.NewServer(svc)

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}

}

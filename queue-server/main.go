package main

import (
	"log"
)

func main() {
	server := NewServer()
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

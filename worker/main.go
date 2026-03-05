package main

import (
	"flag"
	"log"
	"sync"
)

func main() {
	addr := flag.String("queue", "localhost:8080", "queue server address")
	input := flag.String("input", "../files/input.txt", "file to read from")
	output := flag.String("output", "../files/output.txt", "file to write to")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := produce(*addr, *input); err != nil {
			log.Fatalf("producer error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := consume(*addr, *output); err != nil {
			log.Fatalf("consumer error: %v", err)
		}
	}()

	wg.Wait()
	log.Println("done — input and output should be identical")
}

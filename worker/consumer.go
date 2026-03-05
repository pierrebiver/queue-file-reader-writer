package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"queue-file-reader-writer.com/internal/client"
	"queue-file-reader-writer.com/internal/command"
)

func consume(addr, outputPath string) error {
	c, err := client.New(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for {
		line, err := c.Pop()
		if err != nil {
			return err
		}

		if command.EOF.Is(line) {
			log.Println("consumer: received EOF sentinel, done")
			return nil
		}
		if _, err := fmt.Fprintf(writer, "%s\n", line); err != nil {
			return err
		}
		log.Printf("consumer: wrote %q", line)
	}
}

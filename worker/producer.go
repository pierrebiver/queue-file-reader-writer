package main

import (
	"bufio"
	"log"
	"os"

	"queue-file-reader-writer.com/internal/client"
)

func produce(addr, inputPath string) error {
	c, err := client.New(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if err := c.Push(line); err != nil {
			return err
		}
		log.Printf("producer: pushed %q", line)
	}

	return c.SendEOF()
}

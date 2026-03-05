package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"queue-file-reader-writer.com/internal/command"
)

type Client struct {
	conn    net.Conn
	scanner *bufio.Scanner
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to queue server at %s: %w", addr, err)
	}
	return &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
	}, nil
}

func (c *Client) Push(line string) error {
	cmd := command.PushCommand{}
	resp, err := c.send(cmd.CommandName() + " " + line)
	if err != nil {
		return fmt.Errorf("push: %w", err)
	}
	if resp != "OK" {
		return fmt.Errorf("push: unexpected response: %s", resp)
	}
	return nil
}

func (c *Client) Pop() (string, bool, error) {
	cmd := command.PopCommand{}
	resp, err := c.send(cmd.CommandName())
	if err != nil {
		return "", false, fmt.Errorf("pop: %w", err)
	}
	switch {
	case strings.HasPrefix(resp, "DATA "):
		return strings.TrimPrefix(resp, "DATA "), true, nil
	case resp == "EMPTY":
		return "", false, nil
	default:
		return "", false, fmt.Errorf("pop: unexpected response: %s", resp)
	}
}

func (c *Client) SendEOF() error {
	cmd := command.EOFCommand{}
	resp, err := c.send(cmd.CommandName())
	if err != nil {
		return fmt.Errorf("eof: %w", err)
	}
	if resp != "OK" {
		return fmt.Errorf("eof: unexpected response: %s", resp)
	}
	return nil
}

func (c *Client) send(line string) (string, error) {
	if _, err := fmt.Fprintf(c.conn, "%s\n", line); err != nil {
		return "", fmt.Errorf("write error: %w", err)
	}
	if !c.scanner.Scan() {
		return "", fmt.Errorf("no response from server")
	}
	return c.scanner.Text(), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

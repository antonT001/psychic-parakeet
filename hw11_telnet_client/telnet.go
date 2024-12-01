package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrConnNotEstablished = errors.New("connection is not established")

type (
	TelnetClient interface {
		Connect() error
		io.Closer
		Send() error
		Receive() error
	}

	telnetClient struct {
		address string
		timeout time.Duration
		in      io.ReadCloser
		conn    net.Conn
		out     io.Writer
		netS    bufio.Scanner
		localS  bufio.Scanner
	}
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("failed connect to %s", t.address)
	}
	t.conn = conn
	t.netS = *bufio.NewScanner(conn)
	t.localS = *bufio.NewScanner(t.in)
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn == nil {
		return ErrConnNotEstablished
	}
	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return ErrConnNotEstablished
	}

	for t.localS.Scan() {
		_, err := t.conn.Write(append(t.localS.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return ErrConnNotEstablished
	}

	for t.netS.Scan() {
		_, err := t.out.Write(append(t.netS.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}

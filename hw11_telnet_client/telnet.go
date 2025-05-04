package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type Telnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn

	return nil
}

func (t *Telnet) Send() error {
	if _, err := io.Copy(t.conn, t.in); err != nil {
		return err
	}

	return nil
}

func (t *Telnet) Receive() error {
	if _, err := io.Copy(t.out, t.conn); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
	}

	return nil
}

func (t *Telnet) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}

	return nil
}

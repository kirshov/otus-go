package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 30*time.Second, "connect timeout")
	flag.Parse()

	buf := &bytes.Buffer{}
	in := io.NopCloser(buf)

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("address and port is required")
		os.Exit(1)
	}

	tClient := NewTelnetClient(args[0]+":"+args[1], *timeout, in, os.Stdout)
	err := tClient.Connect()
	if err != nil {
		fmt.Printf("Client error: %v\n", err)
		os.Exit(1)
	}

	defer func(tClient TelnetClient) {
		err := tClient.Close()
		if err != nil {
			fmt.Println("Closing telnet client error:", err)
		}
	}(tClient)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go receiveHandler(ctx, &wg, tClient)
	go sendHandler(&wg, tClient, buf)
	wg.Wait()
}

func receiveHandler(ctx context.Context, wg *sync.WaitGroup, t TelnetClient) {
	defer wg.Done()
	go func() {
		<-ctx.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := t.Receive(); err != nil {
				fmt.Println("Receive error:", err)
				return
			}
		}
	}
}

func sendHandler(wg *sync.WaitGroup, t TelnetClient, buf *bytes.Buffer) {
	defer wg.Done()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if sc.Err() != nil {
			return
		}

		buf.WriteString(sc.Text() + "\n")

		if err := t.Send(); err != nil {
			fmt.Println("Send error:", err)
		}
	}
}

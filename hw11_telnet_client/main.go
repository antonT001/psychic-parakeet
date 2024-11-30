package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultTimeout = time.Second * 15
)

var timeout = flag.Duration("timeout", defaultTimeout, "connection timeout")

func main() {
	flag.Parse()
	args := flag.Args()

	address := net.JoinHostPort(args[0], args[1])

	cli := NewTelnetClient(address, *timeout, io.NopCloser(os.Stdin), os.Stdout)
	if err := cli.Connect(); err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh := make(chan os.Signal, 1)
	go func() {
		<-doneCh
		cancel()
	}()

	signal.Notify(doneCh, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cli.Send(); err != nil {
					log.Printf("Send error: %v", err)
					return
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cli.Receive(); err != nil {
					log.Printf("Receive error: %v", err)
					return
				}
			}
		}
	}()

	wg.Wait()
}

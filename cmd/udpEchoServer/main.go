// Copyright 2021 Abhijit Bose. All rights reserved.
// Use of this source code is governed by a Apache 2.0 license that can be found
// in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/boseji/udp"
)

// Synchronization
var wg sync.WaitGroup

func logIt(addr net.Addr, format string, params ...interface{}) {
	s := fmt.Sprintf("%s - ", addr.String())
	log.Printf(s+format, params...)
}

func server(ctx context.Context, u *udp.UDPClient) {
	defer wg.Done()

	log.Println("Server Started on", u.LocalAddr().String())
	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := u.Receive(buf)
			if err != nil {
				// Timeouts are expected
				if strings.Contains(err.Error(), "i/o timeout") {
					continue
				}
				log.Println("Got error in receive - ", err)
				return
			}
			logIt(u.RemoteAddr, "Received %d bytes - %q", n, string(buf[:n]))
			addr, err := net.ResolveUDPAddr("udp", u.RemoteAddr.String())
			if err != nil {
				log.Println("Got error in coverting remote address for UDP -", err)
				return
			}
			n, err = u.Transmit(addr, buf[:n])
			if err != nil {
				log.Println("Got error in transmit - ", err)
				return
			}
			logIt(u.RemoteAddr, "Transmitted %d bytes", n)
		}
	}
}

func main() {
	var port int
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of %s: \n", os.Args[0])
		progName := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "\n# UDP Echo Server Application = '%s'\n\n", progName)
		fmt.Fprint(os.Stderr, "  This program runs a UDP Server on a specified port.\n")
		fmt.Fprint(os.Stderr, "  It echos back the received messages. \n")
		fmt.Fprint(os.Stderr, "  It has maximum buffer size of 2048 bytes.\n\n")
		fmt.Fprint(os.Stderr, "  To terminate the program press 'Ctrl + c' or send SIGINT.\n\n")

		fmt.Fprintf(os.Stderr, "Command: %s", progName)
		flag.VisitAll(func(f *flag.Flag) {
			rg := regexp.MustCompile(`^\*flag\.(.*)Value$`)
			vt := fmt.Sprintf("%T", f.Value)
			m := rg.FindStringSubmatch(vt)
			fmt.Fprintf(os.Stderr, " -%v %v", f.Name, m[1])
		})
		fmt.Fprint(os.Stderr, "\n\n")
		fmt.Fprint(os.Stderr, "Parameters:\n\n")
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n\n")
	}
	flag.IntVar(&port, "p", udp.LocalUDPport, "UDP Local Port range from 1024 to 65535")
	flag.Parse()

	u, err := udp.NewUDPClient(&net.UDPAddr{Port: port})
	if err != nil {
		log.Fatalln("Failed to open Client -", err)
	}
	defer func() {
		u.Close()
		log.Println("UDP Server closed")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	wg.Add(1)
	// Server
	go server(ctx, u)

	// Ctrl+C handler
	go func() {
		select {
		// Block until a signal is received.
		case <-c:
			fmt.Println()
			cancel()
		case <-ctx.Done():
		}
	}()

	// Wait for Everything to Complete
	wg.Wait()
}

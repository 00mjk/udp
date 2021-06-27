// Copyright 2021 Abhijit Bose. All rights reserved.
// Use of this source code is governed by a Apache 2.0 license that can be found
// in the LICENSE file.

package udp

import (
	"context"
	"net"
	"sync"
	"testing"
)

const (

	// maxBufferSize specifies the size of the buffers that
	// are used to temporarily hold data from the UDP packets
	// that we receive.
	maxBufferSize = 1024

	// testingPort defines the test case server port
	testingPort = 8512
)

func TestUDPClient_Create(t *testing.T) {
	u, err := NewUDPClient(nil)
	if err != nil {
		t.Log("failed to create udp client -", err)
		t.Fail()
		return
	}
	u.Close()

	u, err = NewUDPClient(&net.UDPAddr{Port: 62000})
	if err != nil {
		t.Log("failed to create udp client with port -", err)
		t.Fail()
		return
	}
	u.Close()
}

func TestUDPClient_TxRx(t *testing.T) {
	message := "Speak last, Show respect, power and wisdom shall follow"

	// Used to Safely abort the test in case of failure
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// Sync signal for Receive Goroutine
	start := make(chan bool)

	// Wait Group
	var wg sync.WaitGroup

	wg.Add(1)
	// Receiving Side
	go func(t *testing.T) {
		defer wg.Done()
		u, err := NewUDPClient(&net.UDPAddr{Port: testingPort})
		if err != nil {
			t.Log("failed to create udp client -", err)
			t.Fail()
			return
		}
		defer u.Close()
		select {
		case <-start:
			t.Log("Begin Receive")
			buf := make([]byte, maxBufferSize)
			n, err := u.Receive(buf)
			if err != nil {
				t.Log("failed to read udp client -", err)
				t.Fail()
				return
			}
			t.Log("read successfully -", n, "bytes")
			t.Log("Message:", string(buf))
		case <-ctx.Done():
			return
		}
	}(t)

	wg.Add(1)
	// Transmitter
	go func(t *testing.T) {
		defer wg.Done()
		port2 := testingPort + 1
		u2, err := NewUDPClient(&net.UDPAddr{Port: port2})
		if err != nil {
			t.Log("failed to create udp client 2 -", err)
			t.Fail()
			return
		}
		defer u2.Close()
		t.Log("Begin Transmit")
		start <- true // Signal The Receive
		n, err := u2.Transmit(&net.UDPAddr{Port: testingPort}, []byte(message))
		if err != nil {
			t.Log("failed to write udp client 2-", err)
			t.Fail()
			return
		}
		t.Log("wrote successfully -", n, "bytes")
	}(t)

	wg.Wait()
}

func TestUDPClient_Echo(t *testing.T) {
	localAddr := &net.UDPAddr{Port: testingPort}
	u, err := NewUDPClient(localAddr)
	if err != nil {
		t.Log("failed to create udp client -", err)
		t.Fail()
		return
	}
	defer u.Close()

	t.Log("Open UDP client at", u.LocalAddr().String())

	message := []byte("Drop by drop the bucket gets filled")
	n, err := u.Transmit(localAddr, message)
	if err != nil {
		t.Log("failed to write udp client -", err)
		t.Fail()
		return
	}
	t.Log("wrote successfully -", n, "bytes")

	buf := make([]byte, maxBufferSize)
	n, err = u.Receive(buf)
	if err != nil {
		t.Log("failed to read udp client -", err)
		t.Fail()
		return
	}
	t.Log("read successfully -", n, "bytes")
}

func TestUDPClient_Errors(t *testing.T) {
	t.Run("Nil Conn for LocalAddr", func(t *testing.T) {
		u := &UDPClient{}
		a := u.LocalAddr()
		if a != nil {
			t.Errorf("expected nil got %v", a)
		}
	})
	t.Run("Transmit on Nil UDPClient", func(t *testing.T) {
		var u *UDPClient
		_, err := u.Transmit(&net.UDPAddr{Port: testingPort}, []byte("testing"))
		if err == nil {
			t.Error("expected Error got nil")
		}
	})

	t.Run("Transmit on Uninitialized UDPClient", func(t *testing.T) {
		u := &UDPClient{}
		_, err := u.Transmit(&net.UDPAddr{Port: testingPort}, []byte("testing"))
		if err == nil {
			t.Error("expected Error got nil")
		}
	})

	t.Run("Transmit on with wrong inputs UDPClient", func(t *testing.T) {
		u := &UDPClient{}
		_, err := u.Default(nil)
		if err != nil {
			t.Error("failed to initialize got", err)
			return
		}
		defer u.Close()

		_, err = u.Transmit(nil, []byte("testing"))
		if err == nil {
			t.Error("expected Error(missing addr) got nil")
			return
		}

		_, err = u.Transmit(&net.UDPAddr{Port: testingPort}, []byte{})
		if err == nil {
			t.Error("expected Error(missing data) got nil")
			return
		}

	})

	t.Run("Receive on Nil UDPClient", func(t *testing.T) {
		var u *UDPClient
		_, err := u.Receive(make([]byte, maxBufferSize))
		if err == nil {
			t.Error("expected Error got nil")
		}
	})

	t.Run("Receive on Uninitialized UDPClient", func(t *testing.T) {
		u := &UDPClient{}
		_, err := u.Receive(make([]byte, maxBufferSize))
		if err == nil {
			t.Error("expected Error got nil")
		}
	})

	t.Run("Receive on with wrong inputs UDPClient", func(t *testing.T) {
		u := &UDPClient{}
		_, err := u.Default(nil)
		if err != nil {
			t.Error("failed to initialize got", err)
			return
		}
		defer u.Close()

		_, err = u.Receive([]byte{})
		if err == nil {
			t.Error("expected Error got nil")
		}
	})
}

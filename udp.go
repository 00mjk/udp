// Copyright 2021 Abhijit Bose. All rights reserved.
// Use of this source code is governed by a Apache 2.0 license that can be found
// in the LICENSE file.

// Package udp wraps the low level functionality needed for UDP
// Client and Server from `net` package.
package udp

import (
	"fmt"
	"net"
	"time"
)

const (

	// LocalUDPport represents the default local receiving port for UDP client/server
	LocalUDPport = 62048

	// ReadDeadline specifies the UDP default ReadDeadline duration value
	// Note: ReadDeadline is automatically set for every reception.
	ReadDeadline = 50 * time.Millisecond

	// WriteDeadline specifies the UDP default WriteDeadline duration value
	// Note: WriteDeadline is automatically set for every transmission.
	WriteDeadline = 50 * time.Millisecond
)

// UDPClient helps to create a local UDP message sender
// and receiver interface.
type UDPClient struct {
	conn          *net.UDPConn
	ReadDeadline  time.Duration
	WriteDeadline time.Duration
	RemoteAddr    net.Addr
}

// Close helps to close the local UDP client.
// This also implements the io.Closer Interface.
func (u *UDPClient) Close() error {
	defer func() { u.conn = nil }()
	return u.conn.Close()
}

// Default setup the required default values needed for the client to function.
func (u *UDPClient) Default(laddr *net.UDPAddr) (*UDPClient, error) {

	if u == nil {
		u = &UDPClient{
			ReadDeadline:  ReadDeadline,
			WriteDeadline: WriteDeadline,
		}
	}

	if laddr == nil {
		laddr = &net.UDPAddr{Port: LocalUDPport}
	}

	if u.conn == nil {
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			return nil, fmt.Errorf("failed to perform UDP listen in UDPClient - %w", err)
		}
		u.conn = conn
	}

	return u, nil
}

// LocalAddr returns the current local UDP address if the client
// is active. Nil other wise.
func (u *UDPClient) LocalAddr() net.Addr {
	if u != nil && u.conn != nil {
		return u.conn.LocalAddr()
	}
	return nil
}

// Transmit helps to send a block of data to a intended receiver at the specified
// address. This uses the pre-initialized instance of local UDP client.
func (u *UDPClient) Transmit(addr *net.UDPAddr, data []byte) (
	n int,
	err error,
) {
	if u == nil || u.conn == nil {
		err = fmt.Errorf("failed to Transmit due to uninitialized client")
		return
	}

	if addr == nil || len(data) == 0 {
		err = fmt.Errorf("parameter error in Transmit")
		return
	}

	timeout := time.Now().Add(u.WriteDeadline)
	err = u.conn.SetWriteDeadline(timeout)
	if err != nil {
		err = fmt.Errorf("failed in setting write deadline in Transmit - %w", err)
		return
	}

	u.RemoteAddr = addr
	n, err = u.conn.WriteTo(data, addr)
	if err != nil {
		err = fmt.Errorf("failed to write data in Transmit - %w", err)
	}

	return
}

// Receive helps to read data from a remote UDP Server. It also uses the
// open local client for reception
func (u *UDPClient) Receive(rb []byte) (
	n int,
	err error,
) {
	if u == nil || u.conn == nil {
		err = fmt.Errorf("failed to Receive due to uninitialized client")
		return
	}

	if len(rb) == 0 {
		err = fmt.Errorf("parameter error in Receive")
		return
	}

	timeout := time.Now().Add(u.ReadDeadline)
	err = u.conn.SetReadDeadline(timeout)
	if err != nil {
		err = fmt.Errorf("failed in setting read deadline in Receive - %w", err)
		return
	}

	n, addr, err := u.conn.ReadFrom(rb)
	if err != nil {
		err = fmt.Errorf("failed to read data in Receive - %w", err)
	}
	u.RemoteAddr = addr

	return
}

// NewUDPClient creates a local UDP client with a supplied listen port
func NewUDPClient(laddr *net.UDPAddr) (p *UDPClient, err error) {

	// Get the Default values
	p, err = p.Default(laddr)

	return
}

# UDP Client and Server support package

This package wraps the low level functionality needed for UDP
Client and Server from `net` package.

## Usage

> The type `UDPClient` can represent both as a Server or a Client.

### Server ( Receiver )

```go
// This would Open a Server Port at `localhost:62000`
svr, err := udp.NewUDPClient(&net.UDPAddr{Port: 62000})
if err != nil {
    fmt.Prinltn("failed to create udp server -", err)
    return
}
defer svr.Close()

// Receive Data from Client

buf := make([]byte, 1024)   // Buffer to Receive Data
n, err := svr.Receive(buf)  // Waits for max `ReadDeadline` as timeout
if err != nil {
    fmt.Println("failed to read udp client -", err)
    return
}
fmt.Println("read successfully -", n, "bytes")
fmt.Println("Message:", string(buf))
```

> **Note:** `ReadDeadline` can be altered as its part of the `UDPClient`.

### Client ( Transmitter )

```go
// This would Open a Client Port at `localhost:62000`
client, err := udp.NewUDPClient(&net.UDPAddr{Port: 62000})
if err != nil {
    fmt.Prinltn("failed to create udp client -", err)
    return
}
defer client.Close()

// Transmit Data to Server

// Server Address
addr := &net.UDPAddr{IP: net.ParseIP("192.168.1.72"), Port: 62000}
message := "This is the Test message from Client!"
// Waits for max `WriteDeadline` for Transmission to Complete
n, err := client.Transmit(addr, []byte(message)) 
if err != nil {
    fmt.Println("failed to write udp client 2-", err)
    return
}
fmt.Println("wrote successfully -", n, "bytes")

```

> **Note:** `WriteDeadline` can be altered as its part of the `UDPClient`.

## License

```
Copyright 2021 Abhijit Bose. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
// Hello World server in Go
// Binds REP socket to tcp://*:5555
// Expects "Hello" from client, replies with "World"
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"fmt"
	"time"
)

func main() {
	// Socket to talk to clients
	context, _ := zmq.NewContext()
	defer context.Close()

	responder, _ := zmq.NewSocket(zmq.Rep)
	defer responder.Close()

	responder.Bind("tcp://*:5555")

	for {
		// Wait for next request from client
		request, _, _ := responder.RecvPart()
		fmt.Printf("Received request: [%s]\n", string(request))

		// Do some 'work'
		time.Sleep(1 * time.Second)

		// Send reply back to client
		responder.SendPart([]byte("World"), false)
	}
}

// Hello World client in Go
// Connects REQ socket to tcp://localhost:5555
// Sends "Hello" to server, expects "World" back
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"fmt"
)

func main() {

	fmt.Printf("Connecting to hello world server...")

	context, _ := zmq.NewContext()
	defer context.Close()

	requester, _ := context.Socket(zmq.Req)
	defer requester.Close()

	requester.Connect("tcp://localhost:5555")

	for request := 0; request != 10; request++ {

		fmt.Printf("Sending request %d...\n", request)
		requester.SendPart([]byte("Hello"), false)

		reply, _, _ := requester.RecvPart()
		fmt.Printf("Received reply %d: [%s]\n", request, string(reply))
	}
}

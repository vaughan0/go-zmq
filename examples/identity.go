// Demonstrate identities as used by the request-reply pattern. Run this
// program by itself.
package main

import (
	zmq "github.com/vaughan0/go-zmq"
	"github.com/vaughan0/go-zmq/examples/helpers"
)

func main() {
	context, _ := zmq.NewContext()
	defer context.Close()

	sink, _ := context.Socket(zmq.Router)
	defer sink.Close()

	sink.Bind("inproc://example")

	// First allow 0MQ to set the identity
	anonymous, _ := context.Socket(zmq.Req)
	defer anonymous.Close()

	anonymous.Connect("inproc://example")
	anonymous.SendPart([]byte("ROUTER uses a generated UUID"), false)
	helpers.Dump(sink)

	// Then set the identity ourselves
	identified, _ := context.Socket(zmq.Req)
	defer identified.Close()

	identified.SetIdentitiy([]byte("PEER2"))
	identified.Connect("inproc://example")
	identified.SendPart([]byte("ROUTER socket uses REQ's socket identity"), false)
	helpers.Dump(sink)
}

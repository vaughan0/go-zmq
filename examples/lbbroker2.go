// Load-balancing broker
// Demonstrates use of Send and Recv instead of SendPart and RecvPart
package main

import (
	zmq "github.com/vaughan0/go-zmq"
	"github.com/vaughan0/go-zmq/examples/helpers"

	"bytes"
	"time"
)

const (
	nbrClients  = 10
	nbrWorkers  = 3
	workerReady = "\001" // Signals worker is ready
)

func clientTask() {
	context, _ := zmq.NewContext()
	defer context.Close()

	client, _ := context.Socket(zmq.Req)
	defer client.Close()

	client.Connect("ipc://frontend.ipc")

	// Send request, get reply
	for {
		client.SendPart([]byte("HELLO"), false)
		frame, err := client.Recv()
		if err != nil {
			break
		}
		helpers.Dump(frame, "C:")
		time.Sleep(1 * time.Second)
	}
}

// Worker using REQ socket to do load-balancing
func workerTask() {
	context, _ := zmq.NewContext()
	defer context.Close()

	worker, _ := context.Socket(zmq.Req)
	defer worker.Close()

	worker.Connect("ipc://backend.ipc")

	// Tell broker we're ready for work
	worker.SendPart([]byte("READY"), false)

	// Process messages as they arrive
	for {
		frame, err := worker.Recv()
		if err != nil {
			break
		}
		helpers.Dump(frame, "W:")
		frame[2] = []byte("OK")
		worker.Send(frame)
	}
}

// Now we come to the main task. This has the identical functionality to
// the previous {{lbbroker}} broker example
func main() {

	// Prepare our context and sockets
	context, _ := zmq.NewContext()
	defer context.Close()

	frontend, _ := context.Socket(zmq.Router)
	defer frontend.Close()
	frontend.Bind("ipc://frontend.ipc")

	backend, _ := context.Socket(zmq.Router)
	defer backend.Close()
	backend.Bind("ipc://backend.ipc")

	for c := 0; c < nbrClients; c++ {
		go clientTask()
	}

	for w := 0; w < nbrWorkers; w++ {
		go workerTask()
	}

	// Queue of available workers
	workerQueue := make(map[string]bool)

	var poller zmq.PollSet

	poller.Socket(backend, zmq.In)
	poller.Socket(frontend, zmq.In)

	// Here is the main loop for the load balancer. It works the same way
	// as the previous example, but is a lot shorter because it uses Send
	// and Recv methods
	for {
		//  Poll frontend only if we have available workers
		if len(workerQueue) == 0 {
			poller.Update(1, zmq.None)
		} else {
			poller.Update(1, zmq.In)
		}

		_, err := poller.Poll(-1)
		if err != nil {
			break // Interrupted
		}

		// Handle worker activity on backend
		if poller.Events(0).CanRecv() {
			// Use worker identity for load-balancing
			frame, _ := backend.Recv()
			workerId := frame[0]
			clientId := frame[2]
			frame = frame[2:]

			workerQueue[string(workerId)] = true

			// Forward message to client if it's not a READY
			if !bytes.Equal(clientId, []byte("READY")) {
				frontend.Send(frame)
			}
		}

		if poller.Events(1).CanRecv() {
			// Get client request, route to first available worker
			frame, _ := frontend.Recv()

			workerId := dequeu(workerQueue)
			frame = append([][]byte{
				[]byte(workerId),
				[]byte(""),
			}, frame...)
			backend.Send(frame)
		}
	}
}

func dequeu(queue map[string]bool) (id string) {
	for id = range queue {
		delete(queue, id)
		break
	}

	return
}

// Load-balancing broker
// Clients and workers are shown here in-process
package main

import (
	zmq "github.com/vaughan0/go-zmq"
	"github.com/vaughan0/go-zmq/examples/helpers"

	"bytes"
	"fmt"
)

const (
	nbrClients = 10
	nbrWorkers = 3
)

func clientTask() {
	context, _ := zmq.NewContext()
	defer context.Close()

	client, _ := context.Socket(zmq.Req)
	defer client.Close()

	helpers.SetId(client) // Set a printable identity
	client.Connect("ipc://frontend.ipc")

	client.SendPart([]byte("HELLO"), false)
	reply, _, _ := client.RecvPart()
	fmt.Printf("Client: %s\n", reply)
}

// While this example runs in a single process, that is just to make
// it easier to start and stop the example. Each go routine has its own
// context and conceptually acts as a separate process.
// This is the worker task, using a REQ socket to do load-balancing.
func workerTask() {
	context, _ := zmq.NewContext()
	defer context.Close()

	worker, _ := context.Socket(zmq.Req)
	defer worker.Close()

	helpers.SetId(worker) // Set a printable identity
	worker.Connect("ipc://backend.ipc")

	// Tell broker we're ready for work
	worker.SendPart([]byte("READY"), false)

	for {
		// Read and save all frames until we get an empty frame
		// In this example there is only 1, but there could be more
		identity, _, _ := worker.RecvPart()
		worker.RecvPart() // Envelope delimiter

		// Get request, send reply
		request, _, _ := worker.RecvPart()
		fmt.Printf("Worker: %s\n", request)

		worker.SendPart(identity, true)
		worker.SendPart([]byte(""), true)
		worker.SendPart([]byte("OK"), false)
	}
}

// This is the main task. It starts the clients and workers, and then
// routes requests between the two layers. Workers signal READY when
// they start; after that we treat them as ready when they reply with
// a response back to a client. The load-balancing data structure is
// just a queue of next available workers.
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

	var clients int
	for ; clients < nbrClients; clients++ {
		go clientTask()
	}

	for w := 0; w < nbrWorkers; w++ {
		go workerTask()
	}

	// Here is the main loop for the least-recently-used queue. It has two
	// sockets; a frontend for clients and a backend for workers. It polls
	// the backend in all cases, and polls the frontend only when there are
	// one or more workers ready. This is a neat way to use 0MQ's own queues
	// to hold messages we're not ready to process yet. When we get a client
	// reply, we pop the next available worker and send the request to it,
	// including the originating client identity. When a worker replies, we
	// requeue that worker and forward the reply to the original client
	// using the reply envelope.

	// Queue of available workers
	var availableWorkers int
	workerQueue := make(map[string]bool)

	var poller zmq.PollSet

	poller.Socket(backend, zmq.In)
	poller.Socket(frontend, zmq.In)

	for {

		//  Poll frontend only if we have available workers
		if availableWorkers == 0 {
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
			// Queue worker identity for load-balancing
			workerId, _, _ := backend.RecvPart()
			availableWorkers++
			workerQueue[string(workerId)] = true

			// Second frame is empty
			backend.RecvPart()

			// Third frame is READY or else a client reply identity
			clientId, _, _ := backend.RecvPart()

			// If client reply, send rest back to frontend
			if !bytes.Equal(clientId, []byte("READY")) {
				backend.RecvPart() // Empty frame
				reply, _, _ := backend.RecvPart()
				frontend.SendPart(clientId, true)
				frontend.SendPart([]byte(""), true)
				frontend.SendPart(reply, false)

				// Exit after N message
				clients--
				if clients == 0 {
					break
				}
			}
		}

		// Here is how we handle a client request:
		if poller.Events(1).CanRecv() {
			// Now get next client request, route to last-used worker
			// Client request is [identity][empty][request]
			clientId, _, _ := frontend.RecvPart()
			frontend.RecvPart() // Empty frame
			request, _, _ := frontend.RecvPart()

			// Dequeue and drop the next worker identity
			workerId := dequeu(workerQueue)
			availableWorkers--

			backend.SendPart([]byte(workerId), true)
			backend.SendPart([]byte(""), true)
			backend.SendPart(clientId, true)
			backend.SendPart([]byte(""), true)
			backend.SendPart(request, false)
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

// ROUTER-to-DEALER example
package main

import (
	zmq "github.com/vaughan0/go-zmq"
	"github.com/vaughan0/go-zmq/examples/helpers"

	"bytes"
	"fmt"
	"math/rand"
	"time"
)

const nbrWorkers = 10

func workerTask() {

	context, _ := zmq.NewContext()
	defer context.Close()

	worker, _ := context.Socket(zmq.Dealer)
	defer worker.Close()

	// Set a printable identity
	helpers.SetId(worker)
	worker.Connect("tcp://localhost:5671")

	var total int
	for {
		// Tell the broker we're ready for work
		worker.SendPart([]byte(""), true)
		worker.SendPart([]byte("Hi Boss"), false)

		// Get workload from broker, until finished
		worker.RecvPart() // Envelope delimiter
		workload, _, _ := worker.RecvPart()
		finished := bytes.Equal(workload, []byte("Fired!"))
		if finished {
			fmt.Printf("Completed: %d tasks\n", total)
		}
		total++

		// Do some random work
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(500)
		time.Sleep(time.Duration(r) * time.Millisecond)
	}
}

// While this example runs in a single process, that is only to make
// it easier to start and stop the example. Each go routine has its own
// context and conceptually acts as a separate process.
func main() {
	context, _ := zmq.NewContext()
	defer context.Close()

	broker, _ := context.Socket(zmq.Router)
	defer broker.Close()

	broker.Bind("tcp://*:5671")

	for w := 0; w < nbrWorkers; w++ {
		go workerTask()
	}

	// Run for five seconds and then tell workers to end
	endTime := time.Now().Unix() + 5
	workersFired := 0
	for {
		// Next message gives us least recently used worker
		identity, _, _ := broker.RecvPart()
		broker.SendPart(identity, true)
		broker.RecvPart() // Envelope delimiter
		broker.RecvPart() // Response from worker
		broker.SendPart([]byte(""), true)

		// Encourage workers until it's time to fire them
		if time.Now().Unix() < endTime {
			broker.SendPart([]byte("Work harder"), false)
		} else {
			broker.SendPart([]byte("Fired!"), false)
			workersFired++
			if workersFired == nbrWorkers {
				break
			}
		}
	}
}

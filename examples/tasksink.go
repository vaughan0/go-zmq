//
// Task sink
// Binds PULL socket to tcp://localhost:5558
// Collects results from workers via that socket
//
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"fmt"
	"time"
)

func main() {

	// Prepare our context and socket
	context, _ := zmq.NewContext()
	defer context.Close()

	receiver, _ := context.Socket(zmq.Pull)
	defer receiver.Close()
	receiver.Bind("tcp://*:5558")

	// Wait for start of batch
	receiver.Recv()

	// Start our clock now
	startTime := time.Now()

	// Process 100 confirmations
	for taskNbr := 0; taskNbr < 100; taskNbr++ {
		receiver.RecvPart()
		if (taskNbr/10)*10 == taskNbr {
			fmt.Print(":")
		} else {
			fmt.Print(".")
		}
	}

	// Calculate and report duration of batch
	fmt.Printf("Total elapsed time: %d msec\n", time.Now().Sub(startTime)/time.Millisecond)
}

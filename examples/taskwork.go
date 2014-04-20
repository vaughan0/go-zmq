//
// Task worker
// Connects PULL socket to tcp://localhost:5557
// Collects workloads from ventilator via that socket
// Connects PUSH socket to tcp://localhost:5558
// Sends results to sink via that socket
//
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"fmt"
	"strconv"
	"time"
)

func main() {

	context, _ := zmq.NewContext()
	defer context.Close()

	// Socket to receive messages on
	receiver, _ := context.Socket(zmq.Pull)
	defer receiver.Close()

	receiver.Connect("tcp://localhost:5557")

	// Socket to send messages to
	sender, _ := context.Socket(zmq.Push)
	defer sender.Close()

	sender.Connect("tcp://localhost:5558")

	// Process tasks forever
	for {
		str, _, _ := receiver.RecvPart()

		// Simple progress indicator for the viewer
		fmt.Printf("%s.", string(str))

		// Do the work
		duration, _ := strconv.Atoi(string(str))
		time.Sleep(time.Duration(duration) * time.Millisecond)

		// Send results to sink
		sender.SendPart(nil, false)
	}
}

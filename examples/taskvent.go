//
// Task ventilator
// Binds PUSH socket to tcp://localhost:5557
// Sends batch of tasks to workers via that socket
//
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	context, _ := zmq.NewContext()
	defer context.Close()

	// Socket to send messages on
	sender, _ := context.Socket(zmq.Push)
	defer sender.Close()

	sender.Bind("tcp://*:5557")

	// Socket to send start of batch message on
	sink, _ := context.Socket(zmq.Push)
	defer sink.Close()

	sink.Connect("tcp://localhost:5558")

	fmt.Print("Press Enter when the workers are ready; ")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
	fmt.Print("Sending tasks to workers ...")

	// The first message is "0" and signals start of batch
	sink.SendPart([]byte("0"), false)

	// Send 100 tasks
	var taskNbr int
	var totalMsec = 0 // Total expected cost in msec

	for taskNbr = 0; taskNbr < 100; taskNbr++ {
		// Random workload from 1 to 100 msec
		rand.Seed(time.Now().UnixNano())
		workload := rand.Intn(100) + 1
		totalMsec += workload
		str := fmt.Sprintf("%d", workload)
		sender.SendPart([]byte(str), false)
	}

	fmt.Printf("Total expected cost: %d msec\n", totalMsec)
	time.Sleep(time.Second * 1) // Give 0MQ time to deliver
}

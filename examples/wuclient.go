//
// Weather update client
// Connects SUB socket to tcp://localhost:5556
// Collects weather updates and finds avg temp in zipcode
//
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"flag"
	"fmt"
)

func main() {

	context, _ := zmq.NewContext()
	defer context.Close()

	// Socket to talk to server
	fmt.Println("Collecting updates from weather server...")

	subscriber, _ := context.Socket(zmq.Sub)
	defer subscriber.Close()

	subscriber.Connect("tcp://localhost:5556")

	// Subscribe to zipcode, default is NYC, 10001
	filter := "10001 "
	if flag.NArg() > 1 {
		filter = flag.Arg(1)
	}
	subscriber.SetSubscribe([]byte(filter))

	// Process 100 updates
	var updateNbr int64
	var totaltemp int64

	for updateNbr = 0; updateNbr < 100; updateNbr++ {
		str, _, _ := subscriber.RecvPart()

		var zipcode, temperature, relhumidity int64

		fmt.Sscanf(string(str), "%d %d %d", &zipcode, &temperature, &relhumidity)

		totaltemp += temperature
	}

	fmt.Printf("Average temperature for zipcode '%s' was %dF\n", filter, totaltemp/updateNbr)
}

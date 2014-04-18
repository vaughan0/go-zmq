//
// Weather update server
// Binds PUB socket to tcp://*:5556
// Publishes random weather updates
//
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"fmt"
	"math/rand"
)

func main() {

	// Prepare our context and publisher
	context, _ := zmq.NewContext()
	defer context.Close()

	publisher, _ := context.Socket(zmq.Pub)
	defer publisher.Close()

	publisher.Bind("tcp://*:5556")
	//publisher.Bind("ipc://weather.ipc")

	for {
		// Get values that will fool the boss
		zipcode := rand.Intn(100000)
		temperature := rand.Intn(215) - 80
		relhumidity := rand.Intn(50) + 10

		// Send message to all subscribers
		update := fmt.Sprintf("%05d %d %d", zipcode, temperature, relhumidity)
		publisher.SendPart([]byte(update), false)
	}
}

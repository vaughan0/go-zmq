// Report 0MQ version
package main

import (
	zmq "github.com/vaughan0/go-zmq"

	"fmt"
)

func main() {
	major, minor, patch := zmq.Version()
	fmt.Printf("Current 0MQ version is %d.%d.%d\n", major, minor, patch)
}

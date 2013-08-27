package helpers

import (
	zmq "github.com/vaughan0/go-zmq"

	"crypto/rand"
	"fmt"
	"io"
)

// Dump receives all message parts from socket, prints neatly
func Dump(socket *zmq.Socket) {

	fmt.Println("----------------------------------------")

	for {
		// Process parts one by one
		msg, more, err := socket.RecvPart()
		if err != nil {
			fmt.Println(err)
			return
		}

		var binary bool
		for _, x := range msg {
			if x < 32 || x > 127 {
				binary = true
			}
		}

		// Dump the message as text or binary
		fmt.Printf("[%03d] ", len(msg))

		if binary {
			fmt.Printf("% X\n", msg)
		} else {
			fmt.Println(string(msg))
		}

		if !more {
			break // Last message part
		}
	}
}

// SetId sets simple random printable identity on socket
func SetId(socket *zmq.Socket) {
	buf := make([]byte, 4)
	io.ReadFull(rand.Reader, buf)
	id := fmt.Sprintf("%04X-%04X", buf[:2], buf[2:])
	socket.SetIdentitiy([]byte(id))
}

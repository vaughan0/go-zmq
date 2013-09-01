package helpers

import (
	zmq "github.com/vaughan0/go-zmq"

	"crypto/rand"
	"fmt"
	"io"
	"sync"
)

var lock sync.Mutex

// DumpSocket receives all message parts from socket, prints neatly
func DumpSocket(socket *zmq.Socket) {
	frame, err := socket.Recv()
	if err != nil {
		fmt.Println(err)
		return
	}
	Dump(frame, "")
}

// Dump prints frames neatly
func Dump(frame [][]byte, prefix string) {
    lock.Lock()
    defer lock.Unlock()

	fmt.Println("----------------------------------------")

	for _, msg := range frame {
		var binary bool
		for _, x := range msg {
			if x < 32 || x > 127 {
				binary = true
				break
			}
		}

		// Dump the message as text or binary
		fmt.Printf("[%s%03d] ", prefix, len(msg))

		if binary {
			fmt.Printf("% X\n", msg)
		} else {
			fmt.Println(string(msg))
		}
	}
}

// SetId sets simple random printable identity on socket
func SetId(socket *zmq.Socket) (id []byte) {
	buf := make([]byte, 4)
	io.ReadFull(rand.Reader, buf)
	id = []byte(fmt.Sprintf("%04X-%04X", buf[:2], buf[2:]))
	socket.SetIdentitiy(id)

	return
}

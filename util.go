package zmq

import (
	"fmt"
	"sync"
)

// Prefix used for inproc socket pair endpoint names
var PairPrefix = "socket-pair-"

var counter = 0
var lock sync.Mutex

func nextId() int {
	lock.Lock()
	defer lock.Unlock()
	counter++
	return counter
}

func (c *Context) MakePair() (a *Socket, b *Socket, err error) {
	addr := fmt.Sprintf("inproc://%s-%d", PairPrefix, nextId())
	if a, err = c.Socket(Pair); err != nil {
		goto Error
	}
	if err = a.Bind(addr); err != nil {
		goto Error
	}
	if b, err = c.Socket(Pair); err != nil {
		goto Error
	}
	if err = b.Connect(addr); err != nil {
		goto Error
	}
	return

Error:
	if a != nil {
		a.Close()
	}
	if b != nil {
		b.Close()
	}
	return
}

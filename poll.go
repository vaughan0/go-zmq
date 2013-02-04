package zmq

/*
#include <zmq.h>
*/
import "C"

import (
	"time"
)

type Item struct {
	Socket  *Socket
	In      bool
	Out     bool
	CanRecv bool
	CanSend bool
}

func Poll(items []Item, timeout time.Duration) (n int, err error) {

	// Create list of zmq_pollitem_t's
	zitems := make([]C.zmq_pollitem_t, len(items))
	for i, item := range items {
		zitems[i].socket = item.Socket.sock
		if item.In {
			zitems[i].events |= C.ZMQ_POLLIN
		}
		if item.Out {
			zitems[i].events |= C.ZMQ_POLLOUT
		}
	}

	micros := C.long(-1)
	if timeout >= 0 {
		micros = C.long(timeout / time.Microsecond)
	}
	r := C.zmq_poll(&zitems[0], C.int(len(zitems)), micros)
	if r == -1 {
		err = zmqerr()
		return
	}
	n = int(r)

	for i, item := range zitems {
		items[i].CanRecv = (item.revents&C.ZMQ_POLLIN != 0)
		items[i].CanSend = (item.revents&C.ZMQ_POLLOUT != 0)
	}

	return
}

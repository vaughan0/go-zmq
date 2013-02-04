package zmq

/*

#cgo LDFLAGS: -lzmq

#include <zmq.h>
#include <stdlib.h>
#include <string.h>

static int my_errno() {
	return errno;
}

*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

var (
	ErrTerminated = errors.New("zmq context has been terminated")
	ErrTimeout    = errors.New("zmq timeout")
)

type SocketType int

const (
	Req    SocketType = C.ZMQ_REQ
	Rep               = C.ZMQ_REP
	Dealer            = C.ZMQ_DEALER
	Router            = C.ZMQ_ROUTER
	Pub               = C.ZMQ_PUB
	Sub               = C.ZMQ_SUB
	Push              = C.ZMQ_PUSH
	Pull              = C.ZMQ_PULL
	Pair              = C.ZMQ_PAIR
)

/* Context */

type Context struct {
	ctx unsafe.Pointer
}

func NewContextThreads(nthreads int) (ctx *Context, err error) {
	ptr := C.zmq_init(C.int(nthreads))
	if ptr == nil {
		return nil, zmqerr()
	}
	return &Context{ptr}, nil
}

func NewContext() (*Context, error) {
	return NewContextThreads(1)
}

func (c *Context) Close() {
	for {
		r := C.zmq_term(c.ctx)
		if r == -1 {
			if C.my_errno() == C.EINTR {
				continue
			}
			panic(zmqerr())
		}
		break
	}
}

func (c *Context) Socket(socktype SocketType) (sock *Socket, err error) {
	ptr := C.zmq_socket(c.ctx, C.int(socktype))
	if ptr == nil {
		return nil, zmqerr()
	}
	return &Socket{
		ctx:  c,
		sock: ptr,
	}, nil
}

/* Socket */

type Socket struct {
	ctx    *Context
	sock   unsafe.Pointer
	In     <-chan [][]byte
	Out    chan<- [][]byte
	stopch chan bool
	wg     sync.WaitGroup
}

func (s *Socket) Close() {
	if s.stopch != nil {
		close(s.stopch)
		s.wg.Wait()
	}
	C.zmq_close(s.sock)
}

func (s *Socket) Bind(endpoint string) (err error) {
	cstr := C.CString(endpoint)
	defer C.free(unsafe.Pointer(cstr))
	r := C.zmq_bind(s.sock, cstr)
	if r == -1 {
		err = zmqerr()
	}
	return
}
func (s *Socket) Connect(endpoint string) (err error) {
	cstr := C.CString(endpoint)
	defer C.free(unsafe.Pointer(cstr))
	r := C.zmq_connect(s.sock, cstr)
	if r == -1 {
		err = zmqerr()
	}
	return
}

func (s *Socket) SendPart(part []byte, more bool) (err error) {
	var msg C.zmq_msg_t
	toMsg(&msg, part)
	flags := C.int(0)
	if more {
		flags = C.ZMQ_SNDMORE
	}
	r := C.zmq_send(s.sock, &msg, flags)
	if r == -1 {
		err = zmqerr()
	}
	return
}
func (s *Socket) Send(parts [][]byte) (err error) {
	for _, part := range parts[:len(parts)-1] {
		if err = s.SendPart(part, true); err != nil {
			return
		}
	}
	return s.SendPart(parts[len(parts)-1], false)
}

func (s *Socket) RecvPart() (part []byte, more bool, err error) {
	var msg C.zmq_msg_t
	C.zmq_msg_init(&msg)
	r := C.zmq_recv(s.sock, &msg, 0)
	if r == -1 {
		err = zmqerr()
		return
	}
	part = fromMsg(&msg)
	// Check for more parts
	var hasmore C.int64_t
	var optionsz C.size_t = 8
	C.zmq_getsockopt(s.sock, C.ZMQ_RCVMORE, unsafe.Pointer(&hasmore), &optionsz)
	if hasmore != 0 {
		more = true
	}
	return
}
func (s *Socket) Recv() (parts [][]byte, err error) {
	parts = make([][]byte, 0)
	for more := true; more; {
		var part []byte
		if part, more, err = s.RecvPart(); err != nil {
			return
		}
		parts = append(parts, part)
	}
	return
}

/* Utilities */

func zmqerr() error {
	eno := C.my_errno()
	if eno == C.ETERM {
		return ErrTerminated
	} else if eno == C.EAGAIN {
		return ErrTimeout
	}
	str := C.GoString(C.zmq_strerror(eno))
	return errors.New(str)
}

func toMsg(msg *C.zmq_msg_t, data []byte) {
	C.zmq_msg_init_size(msg, C.size_t(len(data)))
	if len(data) > 0 {
		C.memcpy(C.zmq_msg_data(msg), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	}
}
func fromMsg(msg *C.zmq_msg_t) []byte {
	return C.GoBytes(C.zmq_msg_data(msg), C.int(C.zmq_msg_size(msg)))
}

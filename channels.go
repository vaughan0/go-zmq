package zmq

import (
	"log"
)

func (s *Socket) SetupChannels() (err error) {
	if s.stopch != nil {
		panic("SetupChannels() called more than once")
	}

	ctx := s.ctx
	sockin, sockout, err := ctx.MakePair()
	if err != nil {
		return
	}

	// Start worker threads
	s.stopch = make(chan bool)
	s.wg.Add(2)
	incoming := make(chan [][]byte)
	go s.handleIncoming(incoming, sockin)
	outgoing := make(chan [][]byte)
	go s.handleOutgoing(outgoing, sockout)
	s.In = incoming
	s.Out = outgoing

	return
}

func (sock *Socket) handleIncoming(forward chan<- [][]byte, tosend *Socket) {
	defer sock.wg.Done()
	defer tosend.Close()
	defer close(forward)

	var err error
	var parts [][]byte
	items := []Item{
		{
			Socket: sock,
			In:     true,
		}, {
			Socket: tosend,
			In:     true,
		},
	}

	for {
		_, err = Poll(items, -1)
		if err != nil {
			goto Error
		}
		if items[0].CanRecv {
			if parts, err = sock.Recv(); err != nil {
				goto Error
			}
			forward <- parts
		}
		if items[1].CanRecv {
			if parts, err = tosend.Recv(); err != nil {
				goto Error
			}
			if len(parts[0]) > 0 {
				// Close signal
				break
			}
			if err = sock.Send(parts[1:]); err != nil {
				goto Error
			}
		}
	}
	return

Error:
	if err != ErrTerminated {
		log.Println(err)
	}
}

func (s *Socket) handleOutgoing(input <-chan [][]byte, sendto *Socket) {
	defer s.wg.Done()
	defer func() {
		sendto.Send([][]byte{
			{1}, // close signal
		})
		sendto.Close()
	}()

	var err error
	for {
		select {
		case parts := <-input:
			// Send empty part then the rest of the message
			if err = sendto.SendPart([]byte{}, true); err != nil {
				goto Error
			}
			if err = sendto.Send(parts); err != nil {
				goto Error
			}
		case <-s.stopch:
			return
		}
	}
	return

Error:
	if err != ErrTerminated {
		log.Println(err)
	}
}

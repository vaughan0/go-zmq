package zmq

import (
	"testing"
)

var endpoints = []struct {
	endpoint string
	port     uint16
	dyn      bool
}{
	{"tcp://127.0.0.1:5670", 5670, false},
	{"tcp://*:5671", 5671, false},
	{"inproc://my-endpoint", 0, false},
	{"tcp://127.0.0.1:*", 0, true},
	{"tcp://*:*", 0, true},
}

func TestDynBind(t *testing.T) {
	context, err := NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer context.Close()

	socket, err := context.Socket(Push)
	if err != nil {
		t.Fatal(err)
	}
	defer socket.Close()

	err = socket.Bind("tcp://*:*")
	if err != nil {
		t.Fatal(err)
	}

	if socket.Port() != dynPortFrom {
		t.Errorf("expected %d  got %d", dynPortFrom, socket.Port())
	}

	socket2, err := context.Socket(Push)
	if err != nil {
		t.Fatal(err)
	}
	defer socket2.Close()

	err = socket2.Bind("tcp://*:*")
	if err != nil {
		t.Fatal(err)
	}

	if socket2.Port() != dynPortFrom+1 {
		t.Errorf("expected %d  got %d", dynPortFrom+1, socket2.Port())
	}
}

func TestBind(t *testing.T) {

	context, err := NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer context.Close()

	socket, err := context.Socket(Push)
	if err != nil {
		t.Fatal(err)
	}
	defer socket.Close()

	for i, tt := range endpoints {
		var port uint16

		err = socket.Bind(tt.endpoint)
		if err != nil {
			t.Errorf("%d. unexpected err: %s", i, err)
			continue
		}

		port = socket.Port()

		if tt.dyn {
			if port < dynPortFrom || port > dynPortTo {
				t.Errorf("%d. port %d is out of range (%d~%d)", i, port, dynPortFrom, dynPortTo)
			}
			continue
		}

		if port != tt.port {
			t.Errorf("%d. expected %d  got %d", i, tt.port, port)
			continue
		}
	}
}

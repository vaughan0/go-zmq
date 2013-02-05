package zmq

/*
#include <zmq.h>
#include <stdint.h>

#define INT_SIZE (sizeof(int))

*/
import "C"

import (
	"time"
	"unsafe"
)

/* Get */

func (s *Socket) GetType() SocketType {
	return SocketType(s.getInt(C.ZMQ_TYPE))
}
func (s *Socket) GetHWM() uint64 {
	return uint64(s.getInt64(C.ZMQ_HWM))
}
func (s *Socket) GetRecvTimeout() time.Duration {
	return toDuration(int64(s.getInt(C.ZMQ_RCVTIMEO)), time.Millisecond)
}
func (s *Socket) GetSendTimeout() time.Duration {
	return toDuration(int64(s.getInt(C.ZMQ_SNDTIMEO)), time.Millisecond)
}
func (s *Socket) GetAffinity() uint64 {
	return uint64(s.getInt64(C.ZMQ_AFFINITY))
}
func (s *Socket) GetIdentity() []byte {
	return s.getBinary(C.ZMQ_IDENTITY, 256)
}
func (s *Socket) GetRate() (kbits int64) {
	return s.getInt64(C.ZMQ_RATE)
}
func (s *Socket) GetRecoveryIVL() time.Duration {
	return toDuration(s.getInt64(C.ZMQ_RECOVERY_IVL_MSEC), time.Millisecond)
}
func (s *Socket) GetMulticastLoop() bool {
	return s.getInt64(C.ZMQ_MCAST_LOOP) != 0
}
func (s *Socket) GetSendBuffer() (bytes uint64) {
	return uint64(s.getInt64(C.ZMQ_SNDBUF))
}
func (s *Socket) GetRecvBuffer() (bytes uint64) {
	return uint64(s.getInt64(C.ZMQ_RCVBUF))
}
func (s *Socket) GetLinger() time.Duration {
	return toDuration(int64(s.getInt(C.ZMQ_LINGER)), time.Millisecond)
}
func (s *Socket) GetReconnectIVL() time.Duration {
	return toDuration(int64(s.getInt(C.ZMQ_RECONNECT_IVL)), time.Millisecond)
}
func (s *Socket) GetReconnectIVLMax() time.Duration {
	return toDuration(int64(s.getInt(C.ZMQ_RECONNECT_IVL_MAX)), time.Millisecond)
}
func (s *Socket) GetBacklog() int {
	return s.getInt(C.ZMQ_BACKLOG)
}
func (s *Socket) GetFD() uintptr {
	return uintptr(s.getInt(C.ZMQ_FD))
}
func (s *Socket) GetEvents() EventSet {
	return EventSet(s.getInt32(C.ZMQ_EVENTS))
}

/* Set */

func (s *Socket) SetHWM(hwm uint64) {
	s.setInt64(C.ZMQ_HWM, int64(hwm))
}
func (s *Socket) SetAffinity(affinity uint64) {
	s.setInt64(C.ZMQ_AFFINITY, int64(affinity))
}
func (s *Socket) SetIdentitiy(ident []byte) {
	s.setBinary(C.ZMQ_IDENTITY, ident)
}
func (s *Socket) SetRecvTimeout(timeo time.Duration) {
	s.setInt(C.ZMQ_RCVTIMEO, int(fromDuration(timeo, time.Millisecond)))
}
func (s *Socket) SetSendTimeout(timeo time.Duration) {
	s.setInt(C.ZMQ_SNDTIMEO, int(fromDuration(timeo, time.Millisecond)))
}
func (s *Socket) SetRate(kbits int64) {
	s.setInt64(C.ZMQ_RATE, kbits)
}
func (s *Socket) SetRecoveryIVL(ivl time.Duration) {
	s.setInt64(C.ZMQ_RECOVERY_IVL_MSEC, fromDuration(ivl, time.Millisecond))
}
func (s *Socket) SetMulticastLoop(loopback bool) {
	value := int64(0)
	if loopback {
		value = 1
	}
	s.setInt64(C.ZMQ_MCAST_LOOP, value)
}
func (s *Socket) SetSendBuffer(bytes uint64) {
	s.setInt64(C.ZMQ_SNDBUF, int64(bytes))
}
func (s *Socket) SetRecvBuffer(bytes uint64) {
	s.setInt64(C.ZMQ_RCVBUF, int64(bytes))
}
func (s *Socket) SetLinger(linger time.Duration) {
	s.setInt(C.ZMQ_LINGER, int(fromDuration(linger, time.Millisecond)))
}
func (s *Socket) SetReconnectIVL(ivl time.Duration) {
	s.setInt(C.ZMQ_RECONNECT_IVL, int(fromDuration(ivl, time.Millisecond)))
}
func (s *Socket) SetReconnectIVLMax(max time.Duration) {
	s.setInt(C.ZMQ_RECONNECT_IVL_MAX, int(fromDuration(max, time.Millisecond)))
}
func (s *Socket) SetBacklog(backlog int) {
	s.setInt(C.ZMQ_BACKLOG, backlog)
}

/* Utilities */

func (s *Socket) getInt(opt C.int) int {
	var value C.int
	r := C.zmq_getsockopt(s.sock, opt, unsafe.Pointer(&value), nil)
	if r == -1 {
		panic(zmqerr())
	}
	return int(value)
}
func (s *Socket) getInt32(opt C.int) int32 {
	var value C.int32_t
	r := C.zmq_getsockopt(s.sock, opt, unsafe.Pointer(&value), nil)
	if r == -1 {
		panic(zmqerr())
	}
	return int32(value)
}
func (s *Socket) getInt64(opt C.int) int64 {
	var value C.int64_t
	r := C.zmq_getsockopt(s.sock, opt, unsafe.Pointer(&value), nil)
	if r == -1 {
		panic(zmqerr())
	}
	return int64(value)
}
func (s *Socket) getBinary(opt C.int, max int) []byte {
	data := make([]byte, max)
	size := C.size_t(max)
	r := C.zmq_getsockopt(s.sock, opt, unsafe.Pointer(&data[0]), &size)
	if r == -1 {
		panic(zmqerr())
	}
	return data[:int(size)]
}

func (s *Socket) setInt(opt C.int, val int) {
	cval := C.int(val)
	r := C.zmq_setsockopt(s.sock, opt, unsafe.Pointer(&cval), C.INT_SIZE)
	if r == -1 {
		panic(zmqerr())
	}
}
func (s *Socket) setInt64(opt C.int, val int64) {
	r := C.zmq_setsockopt(s.sock, opt, unsafe.Pointer(&val), 8)
	if r == -1 {
		panic(zmqerr())
	}
}
func (s *Socket) setBinary(opt C.int, data []byte) {
	var (
		ptr  unsafe.Pointer
		size C.size_t
	)
	if data != nil && len(data) > 0 {
		ptr = unsafe.Pointer(&data[0])
		size = C.size_t(len(data))
	}
	r := C.zmq_setsockopt(s.sock, opt, ptr, size)
	if r == -1 {
		panic(zmqerr())
	}
}

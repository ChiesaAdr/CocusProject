package socket

import (
	"net"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
const maxBufferSize = 1024

// timeout2resp specifies the timeout to waiting a response(in Seconds)
// const timeout2resp = 4

type Socket struct {
	Conn   *net.UDPConn
	Buffer []byte
}

func NewSocket() *Socket {
	socket := Socket{
		Conn:   &net.UDPConn{},
		Buffer: make([]byte, maxBufferSize),
	}
	return &socket
}

package client

import (
	"cocus/internal/service/socket"
	"cocus/internal/types"
	"io"
	"net"
	"strings"
)

type Client struct {
	socket *socket.Socket
}

func NewClient() *Client {
	return &Client{
		socket: nil,
	}
}

func (s *Client) Run(addr string) {
	//Start socket!
	//  with UDP connection(no guarantee)
	s.socket = socket.NewSocket()

	s.client(addr)
}

//Send control message to the server.
//This message have a pattern accepted only by server
//Example=> Closing10.100.34.12:60716
func (s *Client) SendClose2Server() {

	// Closes the socket after send the notification to server
	defer s.socket.Conn.Close()

	var close []string = []string{("Closing" + s.socket.Conn.LocalAddr().String())}
	//(2)Type Control = control between client and server
	d := types.User{UserType: 2, UserData: close}
	r := strings.NewReader(types.EncodeToString(d))

	io.Copy(s.socket.Conn, r)
}

//Send control message to the server.
//This message have a pattern accepted only by server
//Example=> Delete10.100.34.12:60716=>13
func (s *Client) DeleteMessage2Server(line string) {

	var del []string = []string{("Delete" + s.socket.Conn.LocalAddr().String() + "=>" + line)}
	//(2)Type Control = control between client and server
	d := types.User{UserType: 2, UserData: del}
	r := strings.NewReader(types.EncodeToString(d))
	io.Copy(s.socket.Conn, r)
}

//Client send packet to server, action after the ENTER
func (s *Client) ClientSend(buffer string) {
	var data []string = []string{buffer}
	d := types.User{UserType: 0, UserData: data}
	r := strings.NewReader(types.EncodeToString(d))
	io.Copy(s.socket.Conn, r)
}

//Client wait new packets for the server
func (s *Client) ClientReceive() (user types.User) {
	//Read from UDP Port(Receive packets here)
	nRead, _, err := s.socket.Conn.ReadFrom(s.socket.Buffer)
	if err != nil {
		return
	}
	//Returns the decoded packet for the User type
	return types.DecodeToUser(s.socket.Buffer[:nRead])

}

//Get the actual buffer
func (s *Client) GetBuffer() (buffer []byte) {
	return s.socket.Buffer
}

// client wraps the whole functionality of a UDP client socket
func (s *Client) client(address string) (err error) {

	// Resolve the UDP address so that we can make use of DialUDP
	// with an actual IP and port instead of a address (in case a
	// hostname is specified).
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	// Although we're not in a connection-oriented transport,
	// the act of `dialing` is analogous to the act of performing
	// a `connect(2)` syscall for a socket of type SOCK_DGRAM:
	// - it forces the underlying socket to only read and write
	//   to and from a specific remote address.
	s.socket.Conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	return nil
}

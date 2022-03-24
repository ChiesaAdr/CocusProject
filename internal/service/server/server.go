package server

import (
	"cocus/internal/cache"
	"cocus/internal/config"
	"cocus/internal/service/socket"
	"cocus/internal/types"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Struct to control the client state
type client struct {
	username string
	trans    uint64
	lastTime time.Time
}

type Server struct {
	clients map[string]client
	mu      sync.RWMutex
	cache   *cache.Cache
	socket  *socket.Socket
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]client),
		mu:      sync.RWMutex{},
		cache:   nil,
		socket:  nil,
	}
}

//Basic call to run a new service
func (s *Server) Run(addr string) {
	//Start cache! Using redis
	s.cache = cache.NewCache()

	//Variable to know if you are using a cluster or locally
	if !config.GetScalable() {
		//Monitoring whether Redis is online or not
		go s.cache.RedisMonitor()
	}
	//Start socket!
	//  with UDP connection(no guarantee)
	s.socket = socket.NewSocket()

	go s.serverMonitor()

	//Resolve Address to server
	//Listen UDP socket
	//Loop to wait messages
	go s.server(addr)
}

//Print and Monitoring the activities of each client
func (s *Server) serverMonitor() {
	for {
		s.mu.RLock()
		for v, t := range s.clients {
			log.Printf("Client %v last activity: %s on %v", v, t.username, t.lastTime)
		}
		s.mu.RUnlock()
		time.Sleep(60 * time.Second)
	}
}

//Func to treatment messages of control.
//Delete a specific entry of history
//Delete from list if receive a closing call by client
//Example=> Closing10.100.34.12:60716
func (s *Server) controlClients(msg []string, addr net.Addr) {

	//Control message of type DELETE
	if strings.Contains(msg[0], "Delete"+addr.String()+"=>") {
		re1 := regexp.MustCompile(`\=\>(\d+)`)
		result := re1.FindStringSubmatch(msg[0])

		//Check if find a group!
		if len(result) < 1 {
			return
		}
		//DELETE of cache
		s.cache.RedisDel(addr.String() + result[1])
	}

	//Control message of type CLOSING client
	if strings.Compare(msg[0], "Closing"+addr.String()) == 0 {
		//Needs a Lock here because it delete a entry
		//This MAP is used to control the list of active clients
		s.mu.Lock()
		defer s.mu.Unlock()
		_, exist := s.clients[addr.String()]
		if exist {
			//remove client
			delete(s.clients, addr.String())
		}
		if len(s.clients) == 0 {
			//When all clients disconnect, the DB is flushed.
			s.cache.RedisFlushAll()
		}
	}
}

func (s *Server) updateClients(msg []string, addr net.Addr) {

	//Regex to find username and transtion of each client
	re1 := regexp.MustCompile(`(\w+)\[(\d+)\]:`)
	result := re1.FindStringSubmatch(msg[0])

	//Check if find all group regex
	if len(result) < 3 {
		return
	}
	trans, err := strconv.ParseUint(result[2], 10, 64)
	if err != nil {
		return
	}
	//Need a mutex to WRITE s.clients
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[addr.String()] = client{username: result[1], trans: trans, lastTime: time.Now()}
}

func (s *Server) sendMsg2Clients(msg []string, addr net.Addr) {

	//Regex to find username and transtion of each client
	re1 := regexp.MustCompile(`(\w+)\[(\d+)\]:`)
	result := re1.FindStringSubmatch(msg[0])

	//Check if find all group regex
	if len(result) < 3 {
		return
	}

	//KEYS:
	//addr.String() is Address:Port
	//result[2] is last transation ID
	//s.cache.Order to sort

	//Using the order ID to sort the cache
	//Save message in cache
	s.cache.RedisSet(addr.String()+result[2]+"&"+strconv.FormatUint(s.cache.IncrOrder(), 10), msg[0])
	// s.cache.RedisSet(addr.String()+result[2]+"&"+strconv.FormatUint(s.cache.Order, 10), msg[0])

	//Online User list to send to the Clients
	var users []string
	//Get all the messages of cache(sorted)
	allMsg := s.cache.RedisGetAllValues()

	//READ Lock because it's only read of s.clients
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.clients {
		users = append(users, v.username+"\n")
	}
	//Send messages client by client
	for a := range s.clients {
		udpAddr, _ := net.ResolveUDPAddr("udp", a)
		d := types.User{UserType: 0, UserData: allMsg}
		// Write the packet's contents back to the client.
		s.socket.Conn.WriteToUDP(types.EncodeToBytes(d), udpAddr)
		u := types.User{UserType: 1, UserData: users}
		// Write all online users to the client.
		s.socket.Conn.WriteToUDP(types.EncodeToBytes(u), udpAddr)
	}
}

func (s *Server) readClients() {
	for {
		//Read from UDP Port(Receive packets here)
		n, addr, err := s.socket.Conn.ReadFrom(s.socket.Buffer)
		if err != nil {
			continue
		}
		//Decode []bytes to types.User
		//TODO: Change to use JSON
		user := types.DecodeToUser(s.socket.Buffer[:n])

		//Update client state
		s.updateClients(user.UserData, addr)

		switch user.UserType {
		case 0:
			//(0)Type Message = messages exchanged between clients
			go s.sendMsg2Clients(user.UserData, addr)
		case 2:
			//(2)Type Control = control between client and server
			go s.controlClients(user.UserData, addr)
		default:
		}
	}
}

// server wraps all the UDP echo server functionality.
func (s *Server) server(address string) (err error) {

	log.Printf("Starting Server on %s", address)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}
	s.socket.Conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}

	// `Close`ing the packet "connection" means cleaning the data structures
	// allocated for holding information about the listening Server.
	defer s.socket.Conn.Close()

	//Loop to waiting the messages
	s.readClients()

	return
}

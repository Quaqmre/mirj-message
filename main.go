package main

import (
	"flag"
	"log"
	"net"
)

type Message string

type hub struct {
	clients          map[string]net.Conn
	addClientChan    chan net.Conn
	removeClientChan chan net.Conn
	broadcastChan    chan Message
}

var (
	port = flag.String("port", "9000", "port used for ws connection")
)

func main() {
	flag.Parse()
	hb := newHub()
	lnSCock, err := net.Listen("tcp", ":9001")
	i := 1
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	log.Println("Server running...")
	go hb.run()
	for {
		conn, err := lnSCock.Accept()
		if err != nil {
			log.Fatalln("Error during client connection attemp")
		}
		log.Println("Incoming client connection")
		hb.addClientChan <- conn
		i++
		go func() {
			for {
				recvBuffer := make([]byte, 256)
				bytesRead, err := conn.Read(recvBuffer)
				if err != nil {
					return
				}
				// t := string(m)
				data := recvBuffer[:bytesRead]
				message := Message(data)
				hb.broadcastChan <- message
			}
		}()

	}

}

func (h *hub) run() {
	for {
		select {
		case conn := <-h.addClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.removeClient(conn)
		case m := <-h.broadcastChan:
			h.broadcastMessage(m)
		}
	}
}

func newHub() hub {
	return hub{
		clients:          make(map[string]net.Conn),
		addClientChan:    make(chan net.Conn),
		removeClientChan: make(chan net.Conn),
		broadcastChan:    make(chan Message),
	}
}

func (h *hub) removeClient(conn net.Conn) {
	delete(h.clients, conn.LocalAddr().String())
}

// addClient adds a conn to the pool
func (h *hub) addClient(conn net.Conn) {
	h.clients[conn.RemoteAddr().String()] = conn
}

// broadcastMessage sends a message to all client conns in the pool
func (h *hub) broadcastMessage(m Message) {
	for _, conn := range h.clients {
		b := []byte(m)
		conn.Write(b)
	}
}

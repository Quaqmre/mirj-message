package main

import (
	"log"
	"net"
	"os"

	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/room"
	"github.com/Quaqmre/mırjmessage/user"
)

func main() {
	loggerService := logger.NewLogger(os.Stderr)
	userServce := user.NewUserService(loggerService)
	room := room.NewRoom("deneme", userServce)

	lnSCock, err := net.Listen("tcp", ":9001")
	i := 1
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	log.Println("Server running...")
	go room.Run()
	for {
		conn, err := lnSCock.Accept()
		if err != nil {
			log.Fatalln("Error during client connection attemp")
		}
		log.Println("Incoming client connection")
		room.AddClientChan <- conn
		i++

	}

}

// func (h *hub) run() {
// 	for {
// 		select {
// 		case conn := <-h.addClientChan:
// 			h.addClient(conn)
// 		case conn := <-h.removeClientChan:
// 			h.removeClient(conn)
// 		case m := <-h.broadcastChan:
// 			h.broadcastMessage(m)
// 		}
// 	}
// }

// func newHub() hub {
// 	return hub{
// 		clients:          make(map[string]net.Conn),
// 		addClientChan:    make(chan net.Conn),
// 		removeClientChan: make(chan net.Conn),
// 		broadcastChan:    make(chan Message),
// 	}
// }

// func (h *hub) removeClient(conn net.Conn) {
// 	delete(h.clients, conn.LocalAddr().String())
// }

// // addClient adds a conn to the pool
// func (h *hub) addClient(conn net.Conn) {
// 	h.clients[conn.RemoteAddr().String()] = conn
// }

// // broadcastMessage sends a message to all client conns in the pool
// func (h *hub) broadcastMessage(m Message) {
// 	for _, conn := range h.clients {
// 		b := []byte(m)
// 		conn.Write(b)
// 	}
// }
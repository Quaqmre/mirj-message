package main

import (
	"net/http"
	"os"

	"github.com/Quaqmre/mirjmessage/communication"

	"github.com/Quaqmre/mirjmessage/logger"
	"github.com/Quaqmre/mirjmessage/user"
	"github.com/gorilla/websocket"
)

func main() {
	logger := logger.NewLogger(os.Stderr)
	userService := user.NewUserService(logger)
	server := communication.NewServer(logger, userService)

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Fatal("cmp", "main", "method", "handler", "err", err)
			return
		}
		server.AcceptNewClient(conn)
		// server.Rooms["default"].AddClientChan <- conn
		// server.ac
		logger.Info("cmp", "main", "method", "handler", "msg", "Added new client for default.")
	}

	http.HandleFunc("/", handler)
	logger.Info("cmp", "main", "method", "main", "msg", "Server running..")

	http.ListenAndServe(":9001", nil)

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

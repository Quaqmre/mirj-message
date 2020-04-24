package main

import (
	"github.com/Quaqmre/mırjmessage/communication"
	"log"
	"net/http"
	"os"

	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/user"
	"github.com/gorilla/websocket"
)

func main() {
	loggerService := logger.NewLogger(os.Stderr)
	userServce := user.NewUserService(loggerService)
	rm := communication.NewRoom("deneme", userServce)

	// first handler for each event
	sender := communication.NewSender(rm)
	rm.EventDespatcher.RegisterUserConnectedListener(sender)
	rm.EventDespatcher.RegisterUserLetterListener(sender)

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		rm.AddClientChan <- conn
		log.Println("Added new client. Now", "clients connected.")
	}
	go rm.Run()

	http.HandleFunc("/", handler)
	log.Println("Server running...")

	http.ListenAndServe("localhost:9001", nil)

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

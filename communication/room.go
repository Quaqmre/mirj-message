package communication

import (
	"fmt"
	"sync"

	"github.com/Quaqmre/mırjmessage/events"
	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/pb"
	"github.com/Quaqmre/mırjmessage/user"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// Room is a chating place there is a lot of message and user inside
type Room struct {
	Name             string
	Password         string
	Owners           []int
	Tags             []string
	BanList          []int
	IsPrivite        bool
	Capacity         int
	Clients          map[string]*Client
	AddClientChan    chan *websocket.Conn
	RemoveClientChan chan Client
	// clientService    *client.Service
	userService     user.Service
	logger          logger.Service
	EventDespatcher *EventDispatcher
	mx              sync.RWMutex
	server          *Server
}

// NewRoom give back new chatroom
func NewRoom(name string, u user.Service, logger logger.Service, server *Server) *Room {
	logger.Info("cmp", "room", "method", "NewRoom", "msg", fmt.Sprintf("new room %s Created", name))
	// clientService := client.NewService()
	return &Room{
		Name:             name,
		userService:      u,
		AddClientChan:    make(chan *websocket.Conn, 100),
		RemoveClientChan: make(chan Client),
		EventDespatcher:  NewEventDispatcher(),
		Clients:          make(map[string]*Client),
		logger:           logger,
		server:           server,
	}
}

// Message collapse message
type Package struct {
	UserID  int32  `json:"userid"`
	Message string `json:"message,omitempty"`
}

// Run Handle all messaging issue
func (r *Room) Run() {

	go r.EventDespatcher.RunEventLoop()
	r.logger.Info("cmp", "room", "method", "Run", "msg", fmt.Sprintf("%s room running", r.Name))
	for {
		select {
		case conn := <-r.AddClientChan:
			r.logger.Info("cmp", "room", "method", "Run", "msg", fmt.Sprintf("%s accept a client", r.Name))
			go func() {
				r.acceptNewClient(conn)
			}()

		case c := <-r.RemoveClientChan:
			delete(c.room.Clients, c.ClientIp)
		}
	}
}

func (r *Room) acceptNewClient(conn *websocket.Conn) (err error) {
	defer func() {
		if err != nil {
			r.logger.Fatal("err", fmt.Sprintf("new client cant accept for:%s", err))
		}
	}()

	messageType, data, err := conn.ReadMessage()

	if err != nil {
		return err
	}

	if messageType == websocket.BinaryMessage {

		user, err := r.userService.Marshall(data)
		if err != nil {
			return err
		}
		r.logger.Info("cmp", "room", "method", "acceptNewClient", "msg", fmt.Sprintf("first message is %#v", user))

		newUser, err := r.userService.NewUser(user.Name, user.Password)
		if err != nil {
			return err
		}

		r.logger.Info("cmp", "room", "method", "acceptNewClient", "msg", fmt.Sprintf("new user created:%v", newUser))

		cl, err := NewClient(conn.LocalAddr().String()+string(newUser.UniqID), conn, newUser.UniqID, r, r.server)
		if err != nil {
			return err
		}
		r.mx.Lock()
		r.Clients[cl.ClientIp] = cl
		r.mx.Unlock()
		event := events.UserConnected{ClientID: cl.UserID, Name: newUser.Name}
		r.EventDespatcher.FireUserConnected(&event)
		r.logger.Info("cmp", "room", "method", "acceptNewClient", "msg", fmt.Sprintf("client created:%s", cl.ClientIp))

		cl.Listen()
	}
	return nil
}

// // TODO : channels still not checked
// func (r *Room) handleRead(cl *client.Client) {
// 	log.Println(string(cl.UserID) + " handling Read")
// 	go func() {
// 		ch := make(chan string)
// 		go cl.Listen()
// 		for {
// 			select {
// 			case <-cl.Context.Done():
// 				// delete(r.clientService.Clients, cl.Key)
// 				return
// 			default:
// 				log.Println("waiting coming message from tcp read")
// 				basicMessage := <-ch

// 				username := r.userService.Get(cl.UserID)
// 				msg := fmt.Sprintf("%s: %s", username.Name, basicMessage)
// 				log.Println("formated message created")
// 				r.BroadcastChan <- msg
// 			}
// 		}
// 	}()
// }

func (r *Room) SendToAllClients(message *pb.Message) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	for _, c := range r.Clients {
		c.SendMessage(&bytes)
	}
}
func (r *Room) SendToAllClientsWithIgnored(message *pb.Message, clientIds ...int32) {
	ignored := make(map[int32]bool)

	for _, id := range clientIds {
		ignored[id] = true
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	for _, c := range r.Clients {
		if _, ok := ignored[c.UserID]; !ok {
			c.SendMessage(&bytes)
		}
	}
}
func (r *Room) DeleteClient(key string) {
	r.mx.Lock()
	defer r.mx.Unlock()
	delete(r.Clients, key)
	r.logger.Info("cmp", "room", "method", "DeleteClient", "ms", fmt.Sprintf("client deleted succesfully,key:%s", key))
}

// TODO : uniqId implementasyonuna gerek yoktu çünkü gelen ip uniq

// broadcastMessage sends a message to all client conns in the pool
func (r *Room) broadcastMessage(s string) {
	// log.Println("will send meesage broad cast :" + s)
	// for _, client := range r.Clients {
	// 	err := client.Con.WriteMessage(websocket.BinaryMessage, []byte(s))
	// 	if err != nil {
	// 		log.Println("cant send a client:" + string(client.UserID))
	// 	}
	// }
}

// broadcastMessageWithIgnored sends a message to all client conns in the pool
func (r *Room) broadcastMessageWithIgnored(s string, id ...int32) {
	// log.Println("will send meesage broad cast :" + s)
	// for _, client := range r.Clients {
	// 	for _, i := range id {
	// 		if client.UserID != i {
	// 			err := client.Con.WriteMessage(websocket.BinaryMessage, []byte(s))
	// 			if err != nil {
	// 				r.logger.Fatal(("cant send a client:" + string(client.UserID)))
	// 			}
	// 		}
	// 	}
	// }
}

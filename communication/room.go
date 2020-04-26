package communication

import (
	"fmt"
	"log"

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
	BroadcastChan    chan string
	// clientService    *client.Service
	userService     user.Service
	logger          logger.Service
	EventDespatcher *EventDispatcher
}

// NewRoom give back new chatroom
func NewRoom(name string, u user.Service, logger logger.Service) *Room {
	log.Printf("new room %s Created",name)
	// clientService := client.NewService()
	return &Room{
		Name:             name,
		userService:      u,
		AddClientChan:    make(chan *websocket.Conn, 100),
		RemoveClientChan: make(chan Client),
		BroadcastChan:    make(chan string),
		EventDespatcher:  NewEventDispatcher(),
		Clients:          make(map[string]*Client),
		logger:           logger,
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
	log.Println(r.Name + " room running")
	for {
		select {
		case conn := <-r.AddClientChan:
			log.Println(r.Name + " accept a client")
			go func() {
				r.acceptNewClient(conn)
			}()

		case c := <-r.RemoveClientChan:
			delete(c.server.Clients, c.ClientIp)
		case m := <-r.BroadcastChan:
			log.Println("get broad cast message : ", m)
			r.broadcastMessage(m)
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

		log.Printf("first message is %#v\n", user)

		newUser, err := r.userService.NewUser(user.Name, user.Password)
		if err != nil {
			return err
		}

		log.Printf("new user created:%v\n", newUser)

		cl, err := NewClient(conn.LocalAddr().String()+string(newUser.UniqID), conn, newUser.UniqID, r)
		if err != nil {
			return err
		}
		r.Clients[cl.ClientIp] = cl
		event := events.UserConnected{ClientID: cl.UserID, Name: newUser.Name}
		r.EventDespatcher.FireUserConnected(&event)
		log.Printf("client created:%v\n", cl)

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

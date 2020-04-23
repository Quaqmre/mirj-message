package room

import (
	"log"

	"github.com/Quaqmre/mırjmessage/user"
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
	EventDespatcher *EventDispatcher
}

// NewRoom give back new chatroom
func NewRoom(name string, u user.Service) *Room {
	log.Println("new room Created")
	// clientService := client.NewService()
	return &Room{
		Name:             name,
		userService:      u,
		AddClientChan:    make(chan *websocket.Conn),
		RemoveClientChan: make(chan Client),
		BroadcastChan:    make(chan string),
		EventDespatcher:  NewEventDispatcher(),
		Clients:          make(map[string]*Client),
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
			c.CancelContext()
			_ = c
		case m := <-r.BroadcastChan:
			log.Println("get broad cast message : ", m)
			r.broadcastMessage(m)
		}
	}
}

func (r *Room) acceptNewClient(conn *websocket.Conn) (err error) {
	defer func() {
		if err != nil {
			log.Fatalf("new client cant accept for:%s", err)
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
		event := UserConnected{ClientID: cl.UserID, Name: newUser.Name}
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

// TODO : uniqId implementasyonuna gerek yoktu çünkü gelen ip uniq

// broadcastMessage sends a message to all client conns in the pool
func (r *Room) broadcastMessage(s string) {
	log.Println("will send meesage broad cast :" + s)
	for _, client := range r.Clients {
		err := client.Con.WriteMessage(websocket.BinaryMessage, []byte(s))
		if err != nil {
			log.Fatalln("cant send a client:" + string(client.UserID))
		}
	}
}

package room

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/Quaqmre/mırjmessage/client"
	"github.com/Quaqmre/mırjmessage/user"
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
	AddClientChan    chan net.Conn
	RemoveClientChan chan client.Client
	BroadcastChan    chan string
	clientService    *client.Service
	userService      user.Service
}

// NewRoom give back new chatroom
func NewRoom(name string, u user.Service) *Room {
	log.Println("new room Created")
	clientService := client.NewService()
	return &Room{
		Name:             name,
		clientService:    clientService,
		userService:      u,
		AddClientChan:    make(chan net.Conn),
		RemoveClientChan: make(chan client.Client),
		BroadcastChan:    make(chan string),
	}
}

// Message collapse message
type Package struct {
	UserID  int32  `json:"userid"`
	Message string `json:"message,omitempty"`
}

// Run Handle all messaging issue
func (r *Room) Run() {
	log.Println(r.Name + " room running")
	for {
		select {
		case conn := <-r.AddClientChan:
			log.Println(r.Name + " accept a client")

			recvBuffer := make([]byte, 256)
			bytesRead, err := conn.Read(recvBuffer)
			if err != nil {
				return
			}
			data := recvBuffer[:bytesRead]
			user := &user.User{}
			json.Unmarshal(data, user)
			log.Printf("first message is %#v\n", user)

			newUser, err := r.userService.NewUser(user.Name, user.Password)
			log.Printf("new user created:%v\n", newUser)

			if err != nil {
				return
			}
			cl, _ := r.clientService.New(conn.LocalAddr().String()+string(newUser.UniqID), conn, newUser.UniqID)
			log.Printf("client created:%v\n", cl)
			go func() {

				r.handleRead(cl)
			}()

		case c := <-r.RemoveClientChan:
			close(c.Done)
			_ = c
		case m := <-r.BroadcastChan:
			log.Println("get broad cast message : ", m)
			r.broadcastMessage(m)
		}
	}
}

// TODO : channels still not checked
func (r *Room) handleRead(cl *client.Client) {
	log.Println(string(cl.UserID) + " handleing Read")
	go func() {
		ch := make(chan string)
		go cl.Read(ch)
		for {
			select {
			case <-cl.Done:
				delete(r.clientService.Clients, cl.Key)
				return
			default:
				log.Println("waiting coming message from tcp read")
				basicMessage := <-ch

				username := r.userService.Get(cl.UserID)
				msg := fmt.Sprintf("%s: %s", username.Name, basicMessage)
				log.Println("formated message created")
				r.BroadcastChan <- msg
			}
		}
	}()
}

// TODO : uniqId implementasyonuna gerek yoktu çünkü gelen ip uniq

// broadcastMessage sends a message to all client conns in the pool
func (r *Room) broadcastMessage(s string) {
	log.Println("will send meesage broad cast :" + s)
	for _, client := range r.clientService.Clients {
		_, err := client.Con.Write([]byte(s))
		if err != nil {
			log.Fatalln("cant send a client:" + string(client.UserID))
		}
	}
}

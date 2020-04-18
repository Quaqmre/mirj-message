package room

import (
	"encoding/json"
	"fmt"
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
	userService      *user.UserService
}

// NewRoom give back new chatroom
func NewRoom(name string, u *user.UserService) *Room {
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
	for {
		select {
		case conn := <-r.AddClientChan:
			recvBuffer := make([]byte, 256)
			bytesRead, err := conn.Read(recvBuffer)
			if err != nil {
				return
			}
			data := recvBuffer[:bytesRead]
			user := &user.User{}
			json.Unmarshal(data, user)

			newUser, err := r.userService.NewUser(user.Name, user.Password)

			if err != nil {
				return
			}
			cl, _ := r.clientService.New(conn.LocalAddr().String(), conn, newUser.UniqID)
			go func() {

				r.handleRead(cl)
			}()

		case c := <-r.RemoveClientChan:
			delete(r.clientService.Clients, c.Key)
			_ = c
		case m := <-r.BroadcastChan:
			r.broadcastMessage(m)
		}
	}
}

// TODO : channels still not checked
func (r *Room) handleRead(cl *client.Client) {

	ch := make(chan string)
	go cl.Read(ch)
	for {

		basicMessage := <-ch

		username := r.userService.Dict[cl.UserID]
		msg := fmt.Sprintf("%s: %s", username.Name, basicMessage)
		r.BroadcastChan <- msg
	}
}

// TODO : uniqId implementasyonuna gerek yoktu çünkü gelen ip uniq

// broadcastMessage sends a message to all client conns in the pool
func (r *Room) broadcastMessage(s string) {

	for _, client := range r.clientService.Clients {
		client.Con.Write([]byte(s))
	}
}

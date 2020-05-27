package communication

import (
	"fmt"
	"sync"

	"github.com/Quaqmre/mirjmessage/logger"
	"github.com/Quaqmre/mirjmessage/pb"
	"github.com/Quaqmre/mirjmessage/user"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// Room is a chating place there is a lot of message and user inside
type Room struct {
	Name          string
	Password      string
	Owners        []int
	Tags          []string
	BanList       []int
	IsPrivite     bool
	Capacity      int
	Clients       map[string]*Client
	AddClientChan chan *websocket.Conn
	// clientService    *client.Service
	mx              sync.RWMutex
	server          *Server
	userService     user.Service
	logger          logger.Service
	EventDespatcher EventDispatcher
}

// NewRoom give back new chatroom
func NewRoom(name string, u user.Service, logger logger.Service, server *Server) *Room {
	logger.Info("cmp", "Room", "method", "NewRoom", "msg", fmt.Sprintf("new Room %s Created", name))
	// clientService := client.NewService()
	return &Room{
		Name:            name,
		userService:     u,
		AddClientChan:   make(chan *websocket.Conn, 100),
		EventDespatcher: NewEventDispatcher(),
		Clients:         make(map[string]*Client),
		logger:          logger,
		server:          server,
	}
}

// Run Handle all messaging issue
func (r *Room) Run() {

	go r.EventDespatcher.RunEventLoop()
	r.logger.Info("cmp", "Room", "method", "Run", "msg", fmt.Sprintf("%s Room running", r.Name))
}

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
	r.logger.Info("cmp", "Room", "method", "DeleteClient", "msg", fmt.Sprintf("client deleted succesfully from Room,key:%s", key))
}

func (r *Room) AddClient(client *Client) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.Clients[client.Key] = client
	r.logger.Info("cmp", "Room", "method", "AddClient", "msg", fmt.Sprintf("client added succesfully,key:%s", client.Key))
}

func (r *Room) GetUsers() string {
	users := ""
	if len(r.Clients) > 0 {
		for _, i := range r.Clients {
			users = fmt.Sprintf("%s,%s", i.User.Name, users)

		}
		trimed := users[:len(users)-1]
		users = fmt.Sprintf("%s:%v", trimed, len(r.Clients))

	}
	return users
}

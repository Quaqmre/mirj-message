package communication

import (
	"fmt"
	"sync"

	"github.com/Quaqmre/mirjmessage/events"
	"github.com/Quaqmre/mirjmessage/logger"
	"github.com/Quaqmre/mirjmessage/user"
	"github.com/gorilla/websocket"
)

type Server struct {
	Rooms       map[string]*Room
	Clients     map[string]*Client
	userService user.Service
	logger      logger.Service
	mx          sync.RWMutex
}

func NewServer(logger logger.Service, user user.Service) *Server {
	s := &Server{
		Rooms:       make(map[string]*Room),
		Clients:     make(map[string]*Client),
		userService: user,
		logger:      logger,
	}
	s.CreateRoom("default")
	return s
}
func (s *Server) CreateRoom(name string) *Room {
	rm := NewRoom(name, s.userService, s.logger, s)

	// first handler for each event
	sender := NewSender(rm)
	rm.EventDespatcher.RegisterUserConnectedListener(sender)
	rm.EventDespatcher.RegisterUserLetterListener(sender)
	rm.EventDespatcher.RegisterUserQuitListener(sender)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Rooms[name] = rm
	go rm.Run()

	return rm
}

func (s *Server) GetRooms() string {
	list := ""
	for i := range s.Rooms {
		list = fmt.Sprintf("%s,%s", i, list)
	}
	list = list[:len(list)-1]
	list = fmt.Sprintf("%s:%v", list, len(s.Rooms))

	s.logger.Info("cmp", "server", "method", "GetRooms", "msg", "Rooms listed succesfuly")
	return list
}
func (s *Server) AcceptNewClient(conn *websocket.Conn) (err error) {
	defer func() {
		if err != nil {
			s.logger.Fatal("err", fmt.Sprintf("new client cant accept for:%s", err))
		}
	}()

	messageType, data, err := conn.ReadMessage()

	if err != nil {
		return err
	}

	if messageType == websocket.BinaryMessage {

		user, err := s.userService.Marshall(data)
		if err != nil {
			return err
		}
		s.logger.Info("cmp", "room", "method", "acceptNewClient", "msg", fmt.Sprintf("first message is %#v", user))

		newUser, err := s.userService.NewUser(user.Name, user.Password)
		if err != nil {
			return err
		}

		s.logger.Info("cmp", "room", "method", "acceptNewClient", "msg", fmt.Sprintf("new user created:%v", newUser))
		key := fmt.Sprintf("%s-%v", conn.LocalAddr().String(), newUser.UniqID)
		cl, err := NewClient(key, newUser, conn, newUser.UniqID, s.Rooms["default"], s)
		if err != nil {
			return err
		}
		s.mx.Lock()
		s.Clients[cl.Key] = cl
		s.mx.Unlock()
		event := events.UserConnected{ClientID: cl.UserID, Name: newUser.Name, Key: cl.Key}
		s.Rooms["default"].EventDespatcher.FireUserConnected(&event)
		s.logger.Info("cmp", "room", "method", "acceptNewClient", "msg", fmt.Sprintf("client created:%s", cl.Key))

		cl.Listen()
	}
	return nil
}

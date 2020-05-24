package communication

import (
	"sync"

	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/user"
)

type Server struct {
	Rooms         map[string]*Room
	userService   user.Service
	loggerService logger.Service
	mx            sync.RWMutex
}

func NewServer(logger logger.Service, user user.Service) *Server {
	s := &Server{
		Rooms:         make(map[string]*Room),
		userService:   user,
		loggerService: logger,
	}
	rm := s.CreateRoom("default")
	go rm.Run()
	return s
}
func (s *Server) CreateRoom(name string) *Room {
	rm := NewRoom(name, s.userService, s.loggerService)

	// first handler for each event
	sender := NewSender(rm)
	rm.EventDespatcher.RegisterUserConnectedListener(sender)
	rm.EventDespatcher.RegisterUserLetterListener(sender)
	rm.EventDespatcher.RegisterUserQuitListener(sender)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Rooms[name] = rm
	return rm
}

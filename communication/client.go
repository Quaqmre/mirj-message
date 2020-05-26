package communication

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Quaqmre/mirjmessage/events"
	"github.com/Quaqmre/mirjmessage/logger"
	"github.com/Quaqmre/mirjmessage/pb"
	"github.com/Quaqmre/mirjmessage/user"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var ErrorClientExist = errors.New("user already exist")

// type Service struct {
// 	Clients map[string]*Client
// }

//Client wrap User and add some net info
type Client struct {
	ClientIp      string
	UserID        int32
	User          *user.User
	Con           *websocket.Conn
	Key           string
	ch            chan *[]byte
	Context       context.Context
	cancelContext context.CancelFunc
	room          *Room
	server        *Server
	logger        logger.Service
}

// // NewService make interface of client service
// func NewService() *Service {
// 	return newservice()
// }

// func newservice() *Service {
// 	return &Service{
// 		Clients: make(map[string]*Client),
// 	}
// }

// // New store with user and net connection
// func (c *Service) New(ip string, con *websocket.Conn, userID int32) (*Client, error) {
// 	return c.newClient(ip, con, userID)
// }

// TODO : bir kullanıcı sadece 1 kere mi clients içinde olablir ? Yoksa geçerli olanı mı dönmek gerek
// INFO : client servisi her room özelinde bir tane generete edilmelidir.
func NewClient(key string, user *user.User, con *websocket.Conn, userID int32, room *Room, server *Server) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ClientIp:      con.LocalAddr().String(),
		User:          user,
		Con:           con,
		UserID:        userID,
		Key:           key,
		Context:       ctx,
		cancelContext: cancel,
		room:          room,
		server:        server,
		ch:            make(chan *[]byte, 100),
		logger:        server.logger,
	}

	return client, nil
}

// SendMessage sends game state to the client.
func (c *Client) SendMessage(bytes *[]byte) {
	select {
	case c.ch <- bytes:
	default:
		c.logger.Fatal("err:", "it is dropped message I guess :D")
	}
}

// Done sends done message to the Client which closes the conection.
func (c *Client) Done() {
	c.cancelContext()
}

// // Delete one client in the map
// func (c *Service) Delete(ip string) {
// 	delete(c.Clients, ip)
// }

func (c *Client) Listen() {
	defer c.Con.Close()
	go c.listenWrite()
	c.listenRead()
}

// TODO : channel kapanmalı yoksa hep dinleme yapılacak
func (c *Client) listenRead() {
	for {
		select {
		case <-c.Context.Done():
			c.logger.Warning("cmp", "client", "method", "listenRead", "msg", fmt.Sprintf("%v connectin canceled I cant read any more", c.UserID))
			user := c.server.userService.Get(c.UserID)
			qEvent := events.UserQuit{
				ClientID: c.UserID,
				Name:     user.Name,
				Key:      c.Key,
			}
			c.room.EventDespatcher.FireUserQuit(&qEvent)
			return
		default:
			c.readFromWebSocket()
		}
	}
}

func (c *Client) readFromWebSocket() {
	typ, data, err := c.Con.ReadMessage()

	if err != nil {
		c.logger.Fatal("err:", fmt.Sprintf("when reading message get error from:%v", c.UserID))
		c.cancelContext()
		return
	}

	if typ == websocket.BinaryMessage {
		c.unmarshalUserInput(data)
	}

	// if err != nil {
	// 	log.Fatal("during read message error: ", err)
	// 	return
	// }

}

func (c *Client) unmarshalUserInput(data []byte) {
	protoUserMessage := &pb.UserMessage{}
	if err := proto.Unmarshal(data, protoUserMessage); err != nil {
		c.logger.Fatal("err", fmt.Sprintf("Failed to unmarshal UserInput:%s", err))
		return
	}

	switch x := protoUserMessage.Content.(type) {
	case *pb.UserMessage_Letter:
		userletter := protoUserMessage.GetLetter()
		letterevent := &events.SendLetter{Letter: userletter, ClientId: c.UserID}
		c.room.EventDespatcher.FireUserLetter(letterevent)
	case *pb.UserMessage_Command:
		commandType := protoUserMessage.GetCommand()
		c.handleUserCommand(commandType)
	default:
		log.Fatalf("omg %v", x)
	}
}

// TODO : monitör edebilmek için yazım zamanlarını alıp ortalamasını yazabiliriz.
func (c *Client) listenWrite() {

	c.logger.Info("cmp", "client", "method", "listenWrite", "msg", "Listening write to client")

	for {
		select {
		case bytes := <-c.ch:
			err := c.Con.WriteMessage(websocket.BinaryMessage, *bytes)

			if err != nil {
				//ert.Wrapf(err,fmt.Sprintf("cant send a client:%v" ,c.UserID))
				c.logger.Fatal("err", fmt.Sprintf("cant send a client:%v err:%s", c.UserID, err.Error()))
			}
		case <-c.Context.Done():
			return
		}
	}
}
func (c *Client) handleUserCommand(cmd *pb.Command) {
	switch cmd.Input {
	case pb.Input_LSROOM:
		rooms := c.server.GetRooms()
		letter := &pb.Message_Letter{
			Letter: &pb.Letter{
				Message: rooms,
			},
		}
		message := &pb.Message{Content: letter}

		dat, _ := proto.Marshal(message)
		c.ch <- &dat
	case pb.Input_LSUSER:
		users := c.room.GetUsers()
		letter := &pb.Message_Letter{
			Letter: &pb.Letter{
				Message: users,
			},
		}
		message := &pb.Message{Content: letter}

		dat, _ := proto.Marshal(message)
		c.ch <- &dat
	case pb.Input_CHNAME:
		oldUserName := c.User.Name
		userName := cmd.Message
		err := c.server.userService.ChangeUserName(userName, c.User.UniqID)
		if err != nil {
			c.logger.Warning("cmp", "client", "method", "ChangeUserName", "err", err.Error())
			return
		}
		c.logger.Info("cmp", "client", "method", "ChangeUserName", "msg", fmt.Sprintf("user name %s->%s changed", oldUserName, c.User.Name))
	case pb.Input_JOIN:
		if c.room != nil {
			c.room.EventDespatcher.FireUserQuit(c.QuitEvent())
		}
		roomName := cmd.Message
		room, ok := c.server.Rooms[roomName]
		if !ok {
			fmt.Println(roomName)
			log.Fatal("fatalll")
		}
		c.room = room
		room.EventDespatcher.FireUserConnected(c.ConnectedEvent())

	case pb.Input_MKROOM:
		c.server.CreateRoom(cmd.Message)

	case pb.Input_EXIT:
		c.room.EventDespatcher.FireUserQuit(c.QuitEvent())
		c.room = nil

	}
}

// ConnectedEvent FireUserConnectedEvent
func (c *Client) ConnectedEvent() *events.UserConnected {
	return &events.UserConnected{
		ClientID: c.User.UniqID,
		Key:      c.Key,
		Name:     c.User.Name,
	}
}

// QuitEvent FireUserQuitEvent
func (c *Client) QuitEvent() *events.UserQuit {
	return &events.UserQuit{
		ClientID: c.User.UniqID,
		Key:      c.Key,
		Name:     c.User.Name,
	}
}

// TODO : Muted işlemleri bu katmanda mı handle edilmedi

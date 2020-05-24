package communication

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Quaqmre/mırjmessage/events"
	"github.com/Quaqmre/mırjmessage/pb"
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
	Con           *websocket.Conn
	UserID        int32
	Key           string
	ch            chan *[]byte
	Context       context.Context
	cancelContext context.CancelFunc
	room          *Room
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
func NewClient(ip string, con *websocket.Conn, userID int32, room *Room) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		ClientIp:      ip,
		Con:           con,
		UserID:        userID,
		Key:           con.LocalAddr().String(),
		Context:       ctx,
		cancelContext: cancel,
		room:          room,
		ch:            make(chan *[]byte, 100),
	}

	return client, nil
}

// SendMessage sends game state to the client.
func (c *Client) SendMessage(bytes *[]byte) {
	select {
	case c.ch <- bytes:
	default:
		c.room.logger.Fatal("err:", "it is dropped message I guess :D")
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
			log.Println(string(c.UserID) + "connectin canceled I cant read any more")
			user := c.room.userService.Get(c.UserID)
			qEvent := events.UserQuit{
				ClientID: c.UserID,
				Name:     user.Name,
				Key:      c.ClientIp,
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
		c.room.logger.Fatal("err:", fmt.Sprintf("when reading message get error from:%v", c.UserID))
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
		c.room.logger.Fatal("err", fmt.Sprintf("Failed to unmarshal UserInput:%s", err))
		return
	}

	switch x := protoUserMessage.Content.(type) {
	case *pb.UserMessage_Letter:
		userletter := protoUserMessage.GetLetter()
		letterevent := &events.SendLetter{Letter: userletter, ClientId: c.UserID}
		c.room.EventDespatcher.FireUserLetter(letterevent)
	default:
		log.Fatalf("omg %v", x)
	}
}

// TODO : monitör edebilmek için yazım zamanlarını alıp ortalamasını yazabiliriz.
func (c *Client) listenWrite() {

	log.Println("Listening write to client")

	for {
		select {
		case bytes := <-c.ch:
			err := c.Con.WriteMessage(websocket.BinaryMessage, *bytes)

			if err != nil {
				//ert.Wrapf(err,fmt.Sprintf("cant send a client:%v" ,c.UserID))
				c.room.logger.Fatal("err", fmt.Sprintf("cant send a client:%v err:%s", c.UserID, err.Error()))
			}
		case <-c.Context.Done():
			return
		}
	}
}

// TODO : Muted işlemleri bu katmanda mı handle edilmedi

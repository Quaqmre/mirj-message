package room

import (
	"context"
	"errors"
	"log"

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
	CancelContext context.CancelFunc
	server        *Room
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
		CancelContext: cancel,
		server:        room,
	}

	return client, nil
}

// // Delete one client in the map
// func (c *Service) Delete(ip string) {
// 	delete(c.Clients, ip)
// }

func (c *Client) Listen() {
	defer c.Con.Close()
	// go listenWrite()
	c.listenRead()
}

// TODO : channel kapanmalı yoksa hep dinleme yapılacak
func (c *Client) listenRead() {
	for {
		select {
		case <-c.Context.Done():
			log.Println(string(c.UserID) + "Connectin canceled")
			c.Con.Close()
			return
		default:
			c.readFromWebSocket()
		}
	}
}

func (c *Client) readFromWebSocket() {
	typ, data, err := c.Con.ReadMessage()

	if err != nil {
		c.CancelContext()
		return
	}

	if typ == websocket.BinaryMessage {

		protoUserMessage := &pb.UserMessage{}
		if err := proto.Unmarshal(data, protoUserMessage); err != nil {
			log.Fatalln("Failed to unmarshal UserInput:", err)
			return
		}
		c.server.EventDespatcher.FireUserInput(protoUserMessage)

	}

	// if err != nil {
	// 	log.Fatal("during read message error: ", err)
	// 	return
	// }

}

// TODO : Muted işlemleri bu katmanda mı handle edilmedi

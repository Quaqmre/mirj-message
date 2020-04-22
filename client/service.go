package client

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

var ErrorClientExist = errors.New("user already exist")

type Service struct {
	Clients map[string]*Client
}

//Client wrap User and add some net info
type Client struct {
	ClientIp string
	Con      *websocket.Conn
	UserID   int32
	Key      string
	Done     chan struct{}
	ch       chan *[]byte
}

// NewService make interface of client service
func NewService() *Service {
	return newservice()
}

func newservice() *Service {
	return &Service{
		Clients: make(map[string]*Client),
	}
}

// New store with user and net connection
func (c *Service) New(ip string, con *websocket.Conn, userID int32) (*Client, error) {
	return c.newClient(ip, con, userID)
}

// TODO : bir kullanıcı sadece 1 kere mi clients içinde olablir ? Yoksa geçerli olanı mı dönmek gerek
// INFO : client servisi her room özelinde bir tane generete edilmelidir.
func (c *Service) newClient(ip string, con *websocket.Conn, userID int32) (*Client, error) {

	if _, ok := c.Clients[ip]; ok {
		return nil, ErrorClientExist
	}

	client := &Client{
		ClientIp: ip,
		Con:      con,
		UserID:   userID,
		Key:      con.LocalAddr().String(),
		Done:     make(chan struct{}),
	}

	c.Clients[ip] = client

	return client, nil
}

// Delete one client in the map
func (c *Service) Delete(ip string) {
	delete(c.Clients, ip)
}

// TODO : channel kapanmalı yoksa hep dinleme yapılacak
func (c *Client) Read(ch chan string) {
	for {
		select {
		case <-c.Done:
			log.Println(string(c.UserID) + " done flag expired")
			c.Con.Close()
			return
		default:
			typ, bytesRead, err := c.Con.ReadMessage()

			if err != nil {
				log.Fatal("during read message error: ", err)
				return
			}
			_ = typ
			if typ == websocket.BinaryMessage {
				ch <- string(bytesRead)
			}
		}
	}
}

// TODO : Muted işlemleri bu katmanda mı handle edilmedi

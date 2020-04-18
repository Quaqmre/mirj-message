package client

import (
	"errors"
	"net"
)

var ErrorClientExist = errors.New("user already exist")

type Service struct {
	Clients map[string]*Client
}

//Client wrap User and add some net info
type Client struct {
	ClientIp string
	Con      net.Conn
	UserID   int32
	Key      string
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
func (c *Service) New(ip string, con net.Conn, userID int32) (*Client, error) {
	return c.newClient(ip, con, userID)
}

// TODO : bir kullanıcı sadece 1 kere mi clients içinde olablir ? Yoksa geçerli olanı mı dönmek gerek
// INFO : client servisi her room özelinde bir tane generete edilmelidir.
func (c *Service) newClient(ip string, con net.Conn, userID int32) (*Client, error) {

	if _, ok := c.Clients[ip]; ok {
		return nil, ErrorClientExist
	}

	client := &Client{
		ClientIp: ip,
		Con:      con,
		UserID:   userID,
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
		recvBuffer := make([]byte, 256)
		bytesRead, err := c.Con.Read(recvBuffer)
		if err != nil {
			return
		}
		// t := string(m)
		data := recvBuffer[:bytesRead]

		ch <- string(data)
	}
}

// TODO : Muted işlemleri bu katmanda mı handle edilmedi

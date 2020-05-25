package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Quaqmre/mırjmessage/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var help bool
var host string
var name string
var password string

func init() {
	flag.BoolVar(&help, "help", false, "for command help")
	flag.StringVar(&host, "host", "localhost:9001", "select host")
	flag.StringVar(&name, "name", "anonymous", "select user name")
	flag.StringVar(&password, "pass", "123", "select password")
}

type User struct {
	Name     string
	Password string
}

type Client struct {
	conn *websocket.Conn
}

// This file test for after server up
func main() {
	flag.Parse()
	cl := &Client{}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	u := url.URL{Scheme: "ws", Host: host, Path: "/"}
	fmt.Println(u.String())
	cl.conn, _, _ = websocket.DefaultDialer.Dial(u.String(), nil)
	// if err != nil {
	// log.Fatal("err", err)
	// }
	newU := User{Name: name, Password: password}

	bytes, _ := json.Marshal(newU)

	time.Sleep(time.Second * 1)

	_ = cl.conn.WriteMessage(websocket.BinaryMessage, bytes)
	go func(ch chan os.Signal) {
		_ = <-ch
		fmt.Println("will exit")
		cl.conn.Close()
		os.Exit(0)
		// c.Close()
	}(ch)

	go cl.Listen()
	cl.GetInput()

	time.Sleep(time.Second * 10)

}

func (c *Client) Listen() {
	for {
		select {
		default:
			_, data, erra := c.conn.ReadMessage()
			if erra != nil {
				return
			}
			mes := &pb.Message{}
			err := proto.Unmarshal(data, mes)

			if err != nil {
				log.Println("fatal when un marshal")
			}
			switch mes.Content.(type) {
			case *pb.Message_Letter:
				fmt.Println(mes.GetLetter().Message)
			default:
				fmt.Println(string(data))
				fmt.Println("default")
			}
		}
	}
}

func (c *Client) GetInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		sp := []string{text}
		if text[0] == '&' {
			sp = strings.Split(text, " ")
		}
		switch sp[0] {
		case "&ls":
			if sp[1] == "user" {

				data := c.LSUSER()
				c.MarshalEndWrite(data)
			}
			if sp[1] == "room" {
				data := c.LSROOM()
				c.MarshalEndWrite(data)
			}
		case "&exit":
			data := c.EXIT()
			c.MarshalEndWrite(data)
		case "&ch":
			data := c.CNAME(sp[1])
			c.MarshalEndWrite(data)
		default:
			data := c.Message(text)
			c.MarshalEndWrite(data)
		}
	}
}
func (c *Client) MarshalEndWrite(mes *pb.UserMessage) {
	dat, _ := proto.Marshal(mes)
	_ = c.conn.WriteMessage(websocket.BinaryMessage, dat)
}

func (c *Client) LSROOM() *pb.UserMessage {
	lsroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_LSROOM,
		},
	}
	message := &pb.UserMessage{Content: lsroom}
	return message
}
func (c *Client) LSUSER() *pb.UserMessage {
	lsuser := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_LSUSER,
		},
	}
	message := &pb.UserMessage{Content: lsuser}
	return message
}
func (c *Client) MKROOM(name string) *pb.UserMessage {
	mkroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_MKROOM,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: mkroom}
	return message
}
func (c *Client) RMROOM(name string) *pb.UserMessage {
	rmroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_MKROOM,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: rmroom}
	return message
}
func (c *Client) JOIN(name string) *pb.UserMessage {
	join := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_JOIN,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: join}
	return message
}
func (c *Client) EXIT() *pb.UserMessage {
	exit := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_EXIT,
		},
	}
	message := &pb.UserMessage{Content: exit}
	return message
}
func (c *Client) CNAME(name string) *pb.UserMessage {
	exit := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_CHNAME,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: exit}
	return message
}
func (c *Client) Message(input string) *pb.UserMessage {
	msg := &pb.UserMessage_Letter{
		Letter: &pb.Letter{
			Message: input,
		},
	}
	message := &pb.UserMessage{Content: msg}
	return message
}

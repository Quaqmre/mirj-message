package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Quaqmre/mırjmessage/pb"
	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
	"google.golang.org/protobuf/proto"
)

func (c *Client) Listen(g *gocui.Gui) {
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
				g.Update(func(g *gocui.Gui) error {
					v, err := g.View("messages")
					if err != nil {
						return err
					}
					fmt.Fprint(v, mes.GetLetter().Message)
					return nil
				})
				// fmt.Println(mes.GetLetter().Message)
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

func (c *Client) Message(input string) *pb.UserMessage {
	msg := &pb.UserMessage_Letter{
		Letter: &pb.Letter{
			Message: input,
		},
	}
	message := &pb.UserMessage{Content: msg}
	return message
}

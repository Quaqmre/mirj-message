package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/Quaqmre/mırjmessage/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type User struct {
	Name     string
	Password string
}

// This file test for after server up
func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	u := url.URL{Scheme: "ws", Host: "localhost" + ":9001", Path: "/"}
	fmt.Println(u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("err", err)
	}
	newU := User{Name: "client2", Password: "test"}

	bytes, _ := json.Marshal(newU)

	time.Sleep(time.Second * 5)

	_ = c.WriteMessage(websocket.BinaryMessage, bytes)
	loopbreker := true
	go func(ch chan os.Signal) {
		_ = <-ch
		fmt.Println("will exit")
		c.Close()
		// c.Close()
	}(ch)
	go func() {

		for loopbreker {
			select {
			default:
				_, data, erra := c.ReadMessage()
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
				}
			}
		}
	}()
	lsroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_LSROOM,
		},
	}
	message := &pb.UserMessage{Content: lsroom}

	dat, _ := proto.Marshal(message)
	_ = c.WriteMessage(websocket.BinaryMessage, dat)
	time.Sleep(time.Second * 10)

}

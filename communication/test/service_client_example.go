package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/Quaqmre/mırjmessage/pb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type User struct {
	Name     string
	Password string
}

// This file test for after server up
func main() {
	u := url.URL{Scheme: "ws", Host: "localhost" + ":9001", Path: "/"}
	fmt.Println(u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("err", err.Error())
	}

	newU := User{Name: "akif", Password: "test"}

	bytes, _ := json.Marshal(newU)

	time.Sleep(time.Second * 5)

	_ = c.WriteMessage(websocket.BinaryMessage, bytes)

	userMessage := &pb.UserMessage_Letter{
		Letter: &pb.Letter{
			Message: "selam babalık",
		},
	}
	message := &pb.UserMessage{Content: userMessage}

	datam, err := proto.Marshal(message)
	_ = datam
	// for {

	// 	time.Sleep(time.Second)
	// 	_ = c.WriteMessage(websocket.BinaryMessage, datam)
	// }
	lsroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_LSROOM,
		},
	}
	message = &pb.UserMessage{Content: lsroom}

	dat, _ := proto.Marshal(message)
	for {

		time.Sleep(time.Second)
		_ = c.WriteMessage(websocket.BinaryMessage, dat)
	}
}

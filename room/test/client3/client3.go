package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

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
		log.Fatal(err)
	}
	newU := User{Name: "client3", Password: "test3"}

	bytes, _ := json.Marshal(newU)

	time.Sleep(time.Second * 5)

	_ = c.WriteMessage(websocket.BinaryMessage, bytes)
	for {
		_, message, _ := c.ReadMessage()

		fmt.Println(string(message))
	}
}

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
		log.Fatal("err", err)
	}

	newU := User{Name: "writer", Password: "test"}

	bytes, _ := json.Marshal(newU)

	time.Sleep(time.Second * 5)

	_ = c.WriteMessage(websocket.BinaryMessage, bytes)

	for {

		time.Sleep(time.Second)
		_ = c.WriteMessage(websocket.BinaryMessage, []byte("test"))
	}
}

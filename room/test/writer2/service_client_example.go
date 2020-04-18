package main

import (
	"encoding/json"
	"net"
	"time"
)

type User struct {
	Name     string
	Password string
}

// This file test for after server up
func main() {
	c, err := net.Dial("tcp", "localhost:9001")
	if err != nil {
		panic(err)
	}

	newU := User{Name: "writer2", Password: "writerpas2"}

	bytes, _ := json.Marshal(newU)

	_, _ = c.Write(bytes)

	for {

		time.Sleep(time.Second)
		_, _ = c.Write([]byte("test"))
	}
}

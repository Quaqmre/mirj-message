package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type User struct {
	Name     string
	Password string
}

// This file test for after server up
func main() {
	c, _ := net.Dial("tcp", "localhost:9001")
	newU := User{Name: "client3", Password: "test3"}

	bytes, _ := json.Marshal(newU)

	_, _ = c.Write(bytes)
	for {
		by := make([]byte, 256)
		b, _ := c.Read(by)

		fmt.Println(string(by[:b]))
	}
}

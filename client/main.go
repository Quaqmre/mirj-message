package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
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
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(Layout)
	g.SetViewOnTop("messages")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")
	g.Update(func(g *gocui.Gui) error {
		return nil
	}) // g.SetKeybinding("name", gocui.KeyEnter, gocui.ModNone, Connect)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, cl.Send)

	go cl.Listen(g)
	g.MainLoop()
	// cl.GetInput()

	time.Sleep(time.Second * 10)

}

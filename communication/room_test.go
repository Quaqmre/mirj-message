//Deprecated

package communication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/Quaqmre/mirjmessage/logger"
	"github.com/Quaqmre/mirjmessage/mock"
	"github.com/Quaqmre/mirjmessage/pb"
	"github.com/Quaqmre/mirjmessage/user"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var mockedlogger logger.Service = mock.NewMockedLogger()

var u user.Service = user.NewUserService(mockedlogger)

var roomservice *Room = NewRoom("deneme", u, mockedlogger)

var sender *Sender = NewSender(roomservice)

// Ne kadar kötü bir test case
func TestAtomic_Increase_generate_uniq_Id(t *testing.T) {
	roomservice.EventDespatcher.RegisterUserConnectedListener(sender)
	roomservice.EventDespatcher.RegisterUserLetterListener(sender)
	go roomservice.Run()

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		roomservice.AddClientChan <- conn
		log.Println("Added new client. Now", "clients connected.")
	}

	http.HandleFunc("/", handler)
	log.Println("Server running...")
	var sync sync.WaitGroup
	sync.Add(1)
	go func() {
		http.ListenAndServe("localhost:9001", nil)
	}()
	time.Sleep(time.Second)
	u := url.URL{Scheme: "ws", Host: "localhost" + ":9001", Path: "/"}
	fmt.Println(u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	readedString := make(chan string)
	go func(ch chan string) {

		typ, data, err := c.ReadMessage()
		if err != nil {
			return
		}
		if typ == websocket.BinaryMessage {
			str := ""
			strc := &pb.Message{}
			err := proto.Unmarshal(data, strc)
			if err != nil {
				t.Fatal("cant marshaled bytes")
			}
			switch strc.Content.(type) {
			case *pb.Message_Letter:
				str = strc.GetLetter().Message
			}

			if str != "deneme123" {
				readedString <- str
			}
			readedString <- "success"
		}
	}(readedString)

	newU := &user.User{
		Name:     "akif",
		Password: "deneme",
	}

	bytes, _ := json.Marshal(newU)

	_ = c.WriteMessage(websocket.BinaryMessage, bytes)

	message := &pb.UserMessage{
		Content: &pb.UserMessage_Letter{
			Letter: &pb.Letter{
				Message: "deneme123",
			},
		},
	}

	bytes1, _ := proto.Marshal(message)

	_ = c.WriteMessage(websocket.BinaryMessage, bytes1)

	turned := <-readedString
	sync.Done()

	t.Log("akif")
	if turned != "success" {
		t.Errorf("expected:test: akif but turned:%s", turned)
	}
}

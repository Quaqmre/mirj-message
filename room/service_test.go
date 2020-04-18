package room

import (
	"encoding/json"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/mock"
	"github.com/Quaqmre/mırjmessage/user"
)

var mockedlogger logger.Service = mock.NewMockedLogger()

var u *user.UserService = user.NewUserService(mockedlogger)

var roomservice *Room = NewRoom("deneme", u)

func TestAtomic_Increase_generete_uniq_Id(t *testing.T) {
	readedString := make(chan string)
	go roomservice.Run()
	lnSCock, err := net.Listen("tcp", ":9001")
	if err != nil {
		t.Fatal("not started")
	}

	go func() {
		for {

			conn, err := lnSCock.Accept()
			if err != nil {
				log.Fatalln("Error during client connection attemp")
			}
			roomservice.AddClientChan <- conn
			// log.Println("Incoming client connection")
		}
	}()
	time.Sleep(time.Second * 2)

	c, err := net.Dial("tcp", "localhost:9001")

	if err != nil {
		t.Fatal(err)
	}

	go func(ch chan string) {

		for {

			recvBuffer := make([]byte, 256)
			bytesRead, err := c.Read(recvBuffer)
			if err != nil {
				return
			}
			// t := string(m)
			data := recvBuffer[:bytesRead]
			str := string(data)

			if str != "akif: test" {
				readedString <- str
			}
			readedString <- "success"
		}
	}(readedString)

	newU, _ := roomservice.userService.NewUser("akif", "deneme")

	bytes, _ := json.Marshal(newU)

	_, _ = c.Write(bytes)
	time.Sleep(time.Second)
	_, _ = c.Write([]byte("test"))

	turned := <-readedString

	t.Log("akif")
	if turned != "success" {
		t.Errorf("expected:test: akif but turned:%s", turned)
	}
}

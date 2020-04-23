package room

import (
	"log"

	"github.com/Quaqmre/mırjmessage/pb"
)

type Sender struct {
	server             *Room
	leaderboardCounter int32
}

func NewSender(server *Room) *Sender {
	return &Sender{
		server:             server,
		leaderboardCounter: 0,
	}
}

// HandleUserConnected test
func (sender *Sender) HandleUserConnected(userConnectedEvent *UserConnected) {
	// sender.server.broadcastMessage(fmt.Sprintf("%s-%v user connected server", userConnectedEvent.Name, userConnectedEvent.ClientID))
}

// HandleUserInput test
func (sender *Sender) HandleUserInput(userConnectedEvent *pb.UserMessage) {
	switch x := userConnectedEvent.Content.(type) {
	case *pb.UserMessage_ClientMessage:
		userinput := userConnectedEvent.GetClientMessage()
		sender.server.broadcastMessage(userinput.Message)
	default:
		log.Fatalf("omg %v", x)
	}
}

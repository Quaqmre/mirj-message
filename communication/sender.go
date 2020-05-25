package communication

import (
	"fmt"

	"github.com/Quaqmre/mırjmessage/events"
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
func (sender *Sender) HandleUserConnected(userConnectedEvent *events.UserConnected) {
	text := fmt.Sprintf("%s connected", userConnectedEvent.Name)
	message := &pb.Message{Content: &pb.Message_Letter{Letter: &pb.Letter{Message: text}}}
	sender.server.SendToAllClientsWithIgnored(message)

	// sender.server.broadcastMessage(fmt.Sprintf("%s-%v user connected server", userConnectedEvent.Name, userConnectedEvent.ClientID))
}

// HandleUserInput test
func (sender *Sender) HandleUserLetter(userLettertedEvent *events.SendLetter) {
	user := sender.server.userService.Get(userLettertedEvent.ClientId)
	userLettertedEvent.Letter.Message = fmt.Sprintf("%s:%s", user.Name, userLettertedEvent.Letter.Message)
	message := &pb.Message{Content: &pb.Message_Letter{Letter: userLettertedEvent.Letter}}

	sender.server.SendToAllClientsWithIgnored(message, userLettertedEvent.ClientId)
}

// HandleUserQuit test
func (sender *Sender) HandleUserQuit(userQuitEvent *events.UserQuit) {
	text := fmt.Sprintf("%s quited", userQuitEvent.Name)
	message := &pb.Message{Content: &pb.Message_Letter{Letter: &pb.Letter{Message: text}}}
	sender.server.SendToAllClientsWithIgnored(message, userQuitEvent.ClientID)
	sender.server.DeleteClient(userQuitEvent.Key)

	// sender.server.broadcastMessage(fmt.Sprintf("%s-%v user connected server", userConnectedEvent.Name, userConnectedEvent.ClientID))
}

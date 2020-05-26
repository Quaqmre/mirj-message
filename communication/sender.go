package communication

import (
	"fmt"

	"github.com/Quaqmre/mırjmessage/events"
	"github.com/Quaqmre/mırjmessage/logger"
	"github.com/Quaqmre/mırjmessage/pb"
)

type Sender struct {
	room               *Room
	leaderboardCounter int32
	logger             logger.Service
}

func NewSender(room *Room) *Sender {
	return &Sender{
		room:               room,
		leaderboardCounter: 0,
		logger:             room.logger,
	}
}

// HandleUserConnected test
func (sender *Sender) HandleUserConnected(userConnectedEvent *events.UserConnected) {
	text := fmt.Sprintf("%s connected", userConnectedEvent.Name)
	message := &pb.Message{Content: &pb.Message_Letter{Letter: &pb.Letter{Message: text}}}
	cl := sender.room.server.Clients[userConnectedEvent.Key]
	sender.room.AddClient(cl)

	sender.room.SendToAllClientsWithIgnored(message)
	sender.logger.Info("cmp", "sender", "method", "HandleUserConnected", "msg", "user connected handled succesfully")
	// sender.room.broadcastMessage(fmt.Sprintf("%s-%v user connected room", userConnectedEvent.Name, userConnectedEvent.ClientID))
}

// HandleUserInput test
func (sender *Sender) HandleUserLetter(userLettertedEvent *events.SendLetter) {
	user := sender.room.userService.Get(userLettertedEvent.ClientId)
	userLettertedEvent.Letter.Message = fmt.Sprintf("%s:%s", user.Name, userLettertedEvent.Letter.Message)
	message := &pb.Message{Content: &pb.Message_Letter{Letter: userLettertedEvent.Letter}}
	sender.room.SendToAllClientsWithIgnored(message, userLettertedEvent.ClientId)
	sender.logger.Info("cmp", "sender", "method", "HandleUserLetter", "msg", "user letter handled succesfully")

}

// HandleUserQuit test
func (sender *Sender) HandleUserQuit(userQuitEvent *events.UserQuit) {
	text := fmt.Sprintf("%s quited", userQuitEvent.Name)
	message := &pb.Message{Content: &pb.Message_Letter{Letter: &pb.Letter{Message: text}}}
	sender.room.SendToAllClientsWithIgnored(message, userQuitEvent.ClientID)
	sender.room.DeleteClient(userQuitEvent.Key)
	sender.logger.Info("cmp", "sender", "method", "HandleUserQuit", "msg", "user quit handled succesfully")

	// sender.room.broadcastMessage(fmt.Sprintf("%s-%v user connected room", userConnectedEvent.Name, userConnectedEvent.ClientID))
}

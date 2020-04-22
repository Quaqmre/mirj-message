package room

import (
	"fmt"
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
	sender.server.broadcastMessage(fmt.Sprintf("%s-%v user connected server", userConnectedEvent.Name, userConnectedEvent.ClientID))
}

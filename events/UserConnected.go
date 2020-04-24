package events

import "github.com/Quaqmre/mırjmessage/pb"

type UserConnected struct {
	ClientID int32
	Name     string
}

type SendLetter struct {
	Letter *pb.Letter
	ClientId int32
}
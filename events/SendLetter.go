package events

import "github.com/Quaqmre/mırjmessage/pb"

type SendLetter struct {
	Letter *pb.Letter
	ClientId int32
}
package events

import "github.com/Quaqmre/mirjmessage/pb"

type SendLetter struct {
	Letter   *pb.Letter
	ClientId int32
}

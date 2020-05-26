package main

import "github.com/Quaqmre/mırjmessage/pb"

func (c *Client) LSROOM() *pb.UserMessage {
	lsroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_LSROOM,
		},
	}
	message := &pb.UserMessage{Content: lsroom}
	return message
}
func (c *Client) LSUSER() *pb.UserMessage {
	lsuser := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_LSUSER,
		},
	}
	message := &pb.UserMessage{Content: lsuser}
	return message
}
func (c *Client) MKROOM(name string) *pb.UserMessage {
	mkroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_MKROOM,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: mkroom}
	return message
}
func (c *Client) RMROOM(name string) *pb.UserMessage {
	rmroom := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_MKROOM,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: rmroom}
	return message
}
func (c *Client) JOIN(name string) *pb.UserMessage {
	join := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_JOIN,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: join}
	return message
}
func (c *Client) EXIT() *pb.UserMessage {
	exit := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input: pb.Input_EXIT,
		},
	}
	message := &pb.UserMessage{Content: exit}
	return message
}
func (c *Client) CNAME(name string) *pb.UserMessage {
	exit := &pb.UserMessage_Command{
		Command: &pb.Command{
			Input:   pb.Input_CHNAME,
			Message: name,
		},
	}
	message := &pb.UserMessage{Content: exit}
	return message
}

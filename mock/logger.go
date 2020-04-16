package mock

import "github.com/Quaqmre/mırjmessage/logger"

type mockedlogger struct{}

//NewMockedLogger return nil object
func NewMockedLogger() logger.Service {
	return mockedlogger{}
}

func (m mockedlogger) Info(...interface{}) error {
	return nil
}
func (m mockedlogger) Fatal(...interface{}) error {
	return nil
}
func (m mockedlogger) Warning(...interface{}) error {
	return nil
}

package room

import (
	"time"

	"github.com/Quaqmre/mırjmessage/pb"
)

type UserConnected struct {
	ClientID int32
	Name     string
}

type UserConnectedListener interface {
	HandleUserConnected(*UserConnected)
}

type userConnectedHandler struct {
	event          *UserConnected
	eventListeners []UserConnectedListener
}

func (handler *userConnectedHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserConnected(handler.event)
	}
}

type UserInputListener interface {
	HandleUserInput(*pb.UserMessage)
}

type userInputHandler struct {
	event          *pb.UserMessage
	eventListeners []UserInputListener
}

func (handler *userInputHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserInput(handler.event)
	}
}

type eventHandler interface {
	handle()
}

type EventDispatcher struct {
	running bool

	// EVENT QUEUES

	priority1EventsQueue chan eventHandler

	// LISTENER LISTS

	userConnectedListeners []UserConnectedListener
	userInputListener      []UserInputListener
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		running: false,

		// EVENT QUEUES

		priority1EventsQueue: make(chan eventHandler, 100),

		// LISTENER LISTS

		userConnectedListeners: []UserConnectedListener{},
	}
}
func (dispatcher *EventDispatcher) RunEventLoop() {
	dispatcher.running = true

	for {
		select {

		case handler := <-dispatcher.priority1EventsQueue:
			handler.handle()

		default:
			time.Sleep(time.Millisecond * 1500)
		}
	}
}

func (dispatcher *EventDispatcher) FireUserConnected(event *UserConnected) {
	handler := &userConnectedHandler{
		event:          event,
		eventListeners: dispatcher.userConnectedListeners,
	}

	dispatcher.priority1EventsQueue <- handler
}

func (dispatcher *EventDispatcher) RegisterUserConnectedListener(listener UserConnectedListener) {
	dispatcher.userConnectedListeners = append(dispatcher.userConnectedListeners, listener)
}
func (dispatcher *EventDispatcher) FireUserInput(event *pb.UserMessage) {
	handler := &userInputHandler{
		event:          event,
		eventListeners: dispatcher.userInputListener,
	}

	dispatcher.priority1EventsQueue <- handler
}

func (dispatcher *EventDispatcher) RegisterUserInputListener(listener UserInputListener) {
	dispatcher.userInputListener = append(dispatcher.userInputListener, listener)
}

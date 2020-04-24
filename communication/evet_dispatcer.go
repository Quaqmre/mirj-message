package communication

import (
	"time"

	"github.com/Quaqmre/mırjmessage/events"
)

type UserConnectedListener interface {
	HandleUserConnected(*events.UserConnected)
}

type userConnectedHandler struct {
	event          *events.UserConnected
	eventListeners []UserConnectedListener
}

func (handler *userConnectedHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserConnected(handler.event)
	}
}

type UserLetterListener interface {
	HandleUserLetter(letter *events.SendLetter)
}

type userLetterHandler struct {
	event          *events.SendLetter
	eventListeners []UserLetterListener
}

func (handler *userLetterHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserLetter(handler.event)
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
	userInputListener      []UserLetterListener
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
			time.Sleep(time.Millisecond * 150)
		}
	}
}

func (dispatcher *EventDispatcher) FireUserConnected(event *events.UserConnected) {
	handler := &userConnectedHandler{
		event:          event,
		eventListeners: dispatcher.userConnectedListeners,
	}

	dispatcher.priority1EventsQueue <- handler
}

func (dispatcher *EventDispatcher) RegisterUserConnectedListener(listener UserConnectedListener) {
	dispatcher.userConnectedListeners = append(dispatcher.userConnectedListeners, listener)
}
func (dispatcher *EventDispatcher) FireUserLetter(event *events.SendLetter) {
	handler := &userLetterHandler{
		event:          event,
		eventListeners: dispatcher.userInputListener,
	}

	dispatcher.priority1EventsQueue <- handler
}

func (dispatcher *EventDispatcher) RegisterUserLetterListener(listener UserLetterListener) {
	dispatcher.userInputListener = append(dispatcher.userInputListener, listener)
}

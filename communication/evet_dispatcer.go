package communication

import (
	"time"

	"github.com/Quaqmre/mirjmessage/events"
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
type UserQuitListener interface {
	HandleUserQuit(*events.UserQuit)
}

type userQuitHandler struct {
	event          *events.UserQuit
	eventListeners []UserQuitListener
}

func (handler *userQuitHandler) handle() {
	for _, listener := range handler.eventListeners {
		listener.HandleUserQuit(handler.event)
	}
}

type EventDispatcher struct {
	running bool

	// EVENT QUEUES

	priority1EventsQueue chan eventHandler

	// LISTENER LISTS

	userConnectedListeners []UserConnectedListener
	userInputListener      []UserLetterListener
	userQuitListeners      []UserQuitListener
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		running: false,

		// EVENT QUEUES

		priority1EventsQueue: make(chan eventHandler, 100),

		// LISTENER LISTS

		userConnectedListeners: []UserConnectedListener{},
		userQuitListeners:      []UserQuitListener{},
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

func (dispatcher *EventDispatcher) RegisterUserQuitListener(listener UserQuitListener) {
	dispatcher.userQuitListeners = append(dispatcher.userQuitListeners, listener)
}

func (dispatcher *EventDispatcher) FireUserQuit(event *events.UserQuit) {
	handler := &userQuitHandler{
		event:          event,
		eventListeners: dispatcher.userQuitListeners,
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

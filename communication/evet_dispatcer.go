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

type EventDispatcher interface {
	RunEventLoop()
	FireUserConnected(event *events.UserConnected)
	FireUserLetter(event *events.SendLetter)
	FireUserQuit(event *events.UserQuit)

	RegisterUserConnectedListener(listener UserConnectedListener)
	RegisterUserLetterListener(listener UserLetterListener)
	RegisterUserQuitListener(listener UserQuitListener)
}
type eventDispatcher struct {
	running bool

	// EVENT QUEUES

	priority1EventsQueue chan eventHandler

	// LISTENER LISTS

	userConnectedListeners []UserConnectedListener
	userInputListener      []UserLetterListener
	userQuitListeners      []UserQuitListener
}

// NewEventDispatcher is return interface of dispatcher struct
func NewEventDispatcher() EventDispatcher {
	return &eventDispatcher{
		running: false,

		// EVENT QUEUES

		priority1EventsQueue: make(chan eventHandler, 100),

		// LISTENER LISTS

		userConnectedListeners: []UserConnectedListener{},
		userQuitListeners:      []UserQuitListener{},
	}
}
func (dispatcher *eventDispatcher) RunEventLoop() {
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

func (dispatcher *eventDispatcher) FireUserConnected(event *events.UserConnected) {
	handler := &userConnectedHandler{
		event:          event,
		eventListeners: dispatcher.userConnectedListeners,
	}

	dispatcher.priority1EventsQueue <- handler
}
func (dispatcher *eventDispatcher) RegisterUserConnectedListener(listener UserConnectedListener) {
	dispatcher.userConnectedListeners = append(dispatcher.userConnectedListeners, listener)
}

func (dispatcher *eventDispatcher) RegisterUserQuitListener(listener UserQuitListener) {
	dispatcher.userQuitListeners = append(dispatcher.userQuitListeners, listener)
}

func (dispatcher *eventDispatcher) FireUserQuit(event *events.UserQuit) {
	handler := &userQuitHandler{
		event:          event,
		eventListeners: dispatcher.userQuitListeners,
	}

	dispatcher.priority1EventsQueue <- handler
}

func (dispatcher *eventDispatcher) FireUserLetter(event *events.SendLetter) {
	handler := &userLetterHandler{
		event:          event,
		eventListeners: dispatcher.userInputListener,
	}

	dispatcher.priority1EventsQueue <- handler
}

func (dispatcher *eventDispatcher) RegisterUserLetterListener(listener UserLetterListener) {
	dispatcher.userInputListener = append(dispatcher.userInputListener, listener)
}

type NopEventDespacher struct {
	server *Server
}

func NewNopEventDespacher(server *Server) EventDispatcher {
	return &NopEventDespacher{
		server: server,
	}
}
func (ned *NopEventDespacher) RunEventLoop() {

}
func (ned *NopEventDespacher) FireUserConnected(event *events.UserConnected) {
	ned.server.logger.Info("cmp", "event_dispatcher", "method", "Nop_FireUserConnected", "msg", "Nop dispatcher")
}
func (ned *NopEventDespacher) FireUserLetter(event *events.SendLetter) {
	ned.server.logger.Info("cmp", "event_dispatcher", "method", "Nop_FireUserLetter", "msg", "Nop dispatcher")
}
func (ned *NopEventDespacher) FireUserQuit(event *events.UserQuit) {
	ned.server.logger.Info("cmp", "event_dispatcher", "method", "Nop_FireUserQuit", "msg", "Nop dispatcher")
}
func (ned *NopEventDespacher) RegisterUserConnectedListener(listener UserConnectedListener) {
	ned.server.logger.Info("cmp", "event_dispatcher", "method", "Nop_RegisterUserConnectedListener", "msg", "Nop dispatcher")
}
func (ned *NopEventDespacher) RegisterUserLetterListener(listener UserLetterListener) {
	ned.server.logger.Info("cmp", "event_dispatcher", "method", "Nop_RegisterUserLetterListener", "msg", "Nop dispatcher")
}
func (ned *NopEventDespacher) RegisterUserQuitListener(listener UserQuitListener) {
	ned.server.logger.Info("cmp", "event_dispatcher", "method", "Nop_RegisterUserQuitListener", "msg", "Nop dispatcher")
}

package discordgoself

type EventHandler interface {
	Type() string
	Handle(*Session, interface{})
}

type EventInterfaceProvider interface {
	Type() string
	New() interface{}
}

const interfaceEventType = "__INTERFACE__"
type interfaceEventHandler func(*Session, interface{})

func (eh interfaceEventHandler) Type() string {
	return interfaceEventType
}

func (eh interfaceEventHandler) Handle(s *Session, i interface{}) {
	eh(s, i)
}

var registeredInterfaceProviders = map[string]EventInterfaceProvider{}

func registerInterfaceProvider(eh EventInterfaceProvider) {
	if _, ok := registeredInterfaceProviders[eh.Type()]; ok {
		return
	}
	registeredInterfaceProviders[eh.Type()] = eh
	return
}

type eventHandlerInstance struct {
	eventHandler EventHandler
}

func (s *Session) addEventHandler(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.handlers == nil {
		s.handlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.handlers[eventHandler.Type()] = append(s.handlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

func (s *Session) addEventHandlerOnce(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.onceHandlers == nil {
		s.onceHandlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.onceHandlers[eventHandler.Type()] = append(s.onceHandlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance(eventHandler.Type(), ehi)
	}
}

func (s *Session) AddHandler(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		s.log(LogError, "Invalid handler type, handler will never be called")
		return func() {}
	}

	return s.addEventHandler(eh)
}

func (s *Session) AddHandlerOnce(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		s.log(LogError, "Invalid handler type, handler will never be called")
		return func() {}
	}

	return s.addEventHandlerOnce(eh)
}

func (s *Session) removeEventHandlerInstance(t string, ehi *eventHandlerInstance) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	handlers := s.handlers[t]
	for i := range handlers {
		if handlers[i] == ehi {
			s.handlers[t] = append(handlers[:i], handlers[i+1:]...)
		}
	}

	onceHandlers := s.onceHandlers[t]
	for i := range onceHandlers {
		if onceHandlers[i] == ehi {
			s.onceHandlers[t] = append(onceHandlers[:i], handlers[i+1:]...)
		}
	}
}

func (s *Session) handle(t string, i interface{}) {
	for _, eh := range s.handlers[t] {
		if s.SyncEvents {
			eh.eventHandler.Handle(s, i)
		} else {
			go eh.eventHandler.Handle(s, i)
		}
	}

	if len(s.onceHandlers[t]) > 0 {
		for _, eh := range s.onceHandlers[t] {
			if s.SyncEvents {
				eh.eventHandler.Handle(s, i)
			} else {
				go eh.eventHandler.Handle(s, i)
			}
		}
		s.onceHandlers[t] = nil
	}
}

func (s *Session) handleEvent(t string, i interface{}) {
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()
	s.onInterface(i)
	s.handle(interfaceEventType, i)
	s.handle(t, i)
}

func setGuildIds(g *Guild) {
	for _, c := range g.Channels {
		c.GuildID = g.ID
	}

	for _, m := range g.Members {
		m.GuildID = g.ID
	}

	for _, vs := range g.VoiceStates {
		vs.GuildID = g.ID
	}
}

func (s *Session) onInterface(i interface{}) {
	switch t := i.(type) {
	case *Ready:
		for _, g := range t.Guilds {
			setGuildIds(g)
		}
		s.onReady(t)
	case *GuildCreate:
		setGuildIds(t.Guild)
	case *GuildUpdate:
		setGuildIds(t.Guild)
	case *VoiceServerUpdate:
		go s.onVoiceServerUpdate(t)
	case *VoiceStateUpdate:
		go s.onVoiceStateUpdate(t)
	}
	err := s.State.OnInterface(s, i)
	if err != nil {
		s.log(LogDebug, "error dispatching internal event, %s", err)
	}
}

func (s *Session) onReady(r *Ready) {
	s.sessionID = r.SessionID
}

package session

import (
	"errors"
	"strconv"
	"sync"
)

type Controller struct {
	activeSessionMap map[int64]*Session
	sessionIdCounter int64
	mu               sync.Mutex
}

func NewController() *Controller {
	return &Controller{sessionIdCounter: 1, activeSessionMap: make(map[int64]*Session)}
}

func (c *Controller) NewSession() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	newSess, err := newSession(c.sessionIdCounter)
	if err != nil {
		panic(err)
	}
	c.activeSessionMap[c.sessionIdCounter] = newSess
	c.sessionIdCounter += 1
	return newSess.id
}

func (c *Controller) GetSession(sessionId int64) (*Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	session, ok := c.activeSessionMap[sessionId]
	if ok {
		return session, nil
	}
	return nil, errors.New("Session not found for id = " + strconv.FormatInt(sessionId, 10))
}

func (c *Controller) ProcessCommand(command UserCommand, user *User) error {
	switch command.CommandType {
	case ListSessions:
		user.outboundResponse <- ServerResponse{UserCommandData: command, ErrorMessage: "ListSessions is not implemented yet"}
	case CreateSession:
		sessionId := c.NewSession()
		currentSession, err := c.GetSession(sessionId)
		if err != nil {
			panic("Impossible. Can't find newly created session")
		}
		currentSession.incomingCommands <- SessionCommandCombo{userCommand: command, user: user}
	case JoinSession:
		currentSession, err := c.GetSession(command.Metadata.SessionId)
		if err != nil {
			user.outboundResponse <- ServerResponse{UserCommandData: command, ErrorMessage: "Can't find session with given id"}
			break
		}
		currentSession.incomingCommands <- SessionCommandCombo{userCommand: command, user: user}
	case UnknownCommand:
		user.outboundResponse <- ServerResponse{UserCommandData: command, ErrorMessage: "Unknown command"}
	default:
		if user.session == nil {
			user.outboundResponse <- ServerResponse{UserCommandData: command, ErrorMessage: "User session is not set"}
			break
		}
		user.session.incomingCommands <- SessionCommandCombo{userCommand: command, user: user}
	}
	return nil
}

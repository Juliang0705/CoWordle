package session

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

const (
	RowSize = 6
	ColSize = 5
)

type GuessType string

const (
	NotGuessed    GuessType = "NotGuessed"
	MisPositioned GuessType = "MisPositioned"
	Correct       GuessType = "Correct"
	Incorrect     GuessType = "Incorrect"
)

type Cell struct {
	Value     string
	GuessType GuessType
}

type GameState struct {
	SessionId         int64
	Word              string
	Grid              [RowSize][ColSize]Cell
	CurrentCharacters string
	GuessCount        int
	GameOver          bool
	CurrentUsers      map[string]bool
}

type SessionCommandCombo struct {
	userCommand UserCommand
	user        *User
}

type Session struct {
	id               int64
	gameState        GameState
	incomingCommands chan SessionCommandCombo
	currentUsers     map[string]*User
	unregisterUser   chan *User
}

func newSession(id int64) (*Session, error) {
	session := &Session{
		id:               id,
		currentUsers:     make(map[string]*User),
		gameState:        GameState{SessionId: id, Word: GetWord(), CurrentUsers: make(map[string]bool)},
		incomingCommands: make(chan SessionCommandCombo, 1024),
		unregisterUser:   make(chan *User, 1024)}
	go session.DoSessionLoop()
	return session, nil
}

func (s *Session) HandleCommand(command SessionCommandCombo) error {
	var metadata CommandMetadata = command.userCommand.Metadata
	var user *User = command.user
	if command.userCommand.CommandType != CreateSession && command.userCommand.CommandType != JoinSession {
		if metadata.SessionId == 0 {
			user.outboundResponse <- ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: "Invalid session id"}
			return nil
		}
		if _, ok := s.currentUsers[metadata.UserId]; !ok {
			user.outboundResponse <- ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: "User id not found"}
			return nil
		}
	}
	switch command.userCommand.CommandType {
	case CreateSession, JoinSession:
		s.currentUsers[metadata.UserId] = user
		s.gameState.CurrentUsers[metadata.UserId] = true
		user.session = s
		s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState})
	case AddCharacter:
		if s.gameState.GameOver {
			s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: "Game over. No more guesses"})
			break
		}
		newChar := command.userCommand.AddCharacterCommandPayloadData.Character
		if len(newChar) != 1 {
			s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: "Can only add one character at a time"})
			break
		}
		size := len(s.gameState.CurrentCharacters)
		if size == ColSize {
			s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: "Current character set is full"})
			break
		}
		s.gameState.CurrentCharacters += newChar
		s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState})
	case BackspaceCharacter:
		size := len(s.gameState.CurrentCharacters)
		if size > 0 {
			s.gameState.CurrentCharacters = s.gameState.CurrentCharacters[:size-1]
		}
		s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState})
	case CommitGuess:
		err := s.guess(s.gameState.CurrentCharacters)
		s.gameState.CurrentCharacters = ""
		if err != nil {
			s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: err.Error()})
		} else {
			s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState})
		}
	case ResetSession:
		s.gameState = GameState{SessionId: s.gameState.SessionId, CurrentUsers: s.gameState.CurrentUsers, Word: GetWord()}
		s.BroadcastResponse(ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState})
	default:
		user.outboundResponse <- ServerResponse{UserCommandData: command.userCommand, GameStateData: s.gameState, ErrorMessage: "Unknown command"}
	}
	return nil
}

func (s *Session) BroadcastResponse(response ServerResponse) {
	for _, v := range s.currentUsers {
		v.outboundResponse <- response
	}
}

func (s *Session) DoSessionLoop() {
	for {
		select {
		case user := <-s.unregisterUser:
			var username string
			for k, v := range s.currentUsers {
				if v == user {
					username = k
					break
				}
			}
			delete(s.currentUsers, username)
			delete(s.gameState.CurrentUsers, username)
		case command := <-s.incomingCommands:
			err := s.HandleCommand(command)
			if err != nil {
				log.Printf("error: %v", err)
			}
		}
	}
}

func (s *Session) guess(word string) error {
	word = strings.ToUpper(word)
	if s.gameState.GameOver {
		return errors.New("Game over. No more guesses")
	}
	if s.gameState.GuessCount == RowSize {
		return errors.New("Game over. Too many guesses")
	}
	if len(word) != ColSize {
		return errors.New("Size of word has to be equal to " + strconv.Itoa(ColSize))
	}
	correctCount := 0
	for index := 0; index < ColSize; index += 1 {
		var guess *Cell = &s.gameState.Grid[s.gameState.GuessCount][index]
		guess.Value = string(word[index])
		if word[index] == s.gameState.Word[index] {
			guess.GuessType = Correct
			correctCount += 1
		} else if strings.Contains(s.gameState.Word, string(word[index])) {
			guess.GuessType = MisPositioned
		} else {
			guess.GuessType = Incorrect
		}
	}
	if correctCount == ColSize || s.gameState.GuessCount+1 == RowSize {
		s.gameState.GameOver = true
	}
	s.gameState.GuessCount += 1
	return nil
}

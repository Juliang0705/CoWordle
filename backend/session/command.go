package session

type UserCommandType string

const (
	UnknownCommand     UserCommandType = "UnknownCommand"
	ListSessions       UserCommandType = "ListSessions"
	CreateSession      UserCommandType = "CreateSession"
	JoinSession        UserCommandType = "JoinSession"
	AddCharacter       UserCommandType = "AddCharacter"
	BackspaceCharacter UserCommandType = "BackspaceCharacter"
	CommitGuess        UserCommandType = "CommitGuess"
	ResetSession       UserCommandType = "ResetSession"
)

type ListSessionsCommandPayload struct {
}

type CreateSessionCommandPayload struct {
}

type JoinSessionCommandPayload struct {
}

type AddCharacterCommandPayload struct {
	Character string
}

type BackspaceCharacterCommandPayload struct {
}

type CommitGuessCommandPayload struct {
}

type ResetSessionPayload struct {
}

type CommandMetadata struct {
	SessionId int64
	UserId    string
}

type UserCommand struct {
	CommandType                          UserCommandType
	Metadata                             CommandMetadata
	ListSessionsCommandPayloadData       ListSessionsCommandPayload
	CreateSessionCommandPayLoadData      CreateSessionCommandPayload
	JoinSessionCommandPayloadData        JoinSessionCommandPayload
	AddCharacterCommandPayloadData       AddCharacterCommandPayload
	BackspaceCharacterCommandPayloadData BackspaceCharacterCommandPayload
	CommitGuessCommandPayloadData        CommitGuessCommandPayload
	ResetSessionPayloadData              ResetSessionPayload
}

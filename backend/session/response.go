package session

type ServerResponse struct {
	UserCommandData UserCommand
	GameStateData   GameState
	ErrorMessage    string
}

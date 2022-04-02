export interface CommandMetadata {
	SessionId?: number;
	UserId: string;
}

export interface  AddCharacterCommandPayload  {
	Character: string;
}

export interface UserCommand {
    CommandType: string;
    Metadata: CommandMetadata;
    AddCharacterCommandPayloadData?: AddCharacterCommandPayload;
}

export interface Cell {
    Value: string;
    GuessType: string;
}

export interface GameState {
    SessionId: number;
    Word: string;
    Grid: Cell[][];
    CurrentCharacters: string;
    GuessCount: number;
    GameOver: boolean;
    CurrentUsers: { [key: string]: boolean };
}

export interface ServerResponse {
    UserCommandData: UserCommand;
    GameStateData: GameState;
    ErrorMessage: string;
}
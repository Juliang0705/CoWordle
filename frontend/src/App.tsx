import React, { ChangeEvent } from "react";
import "./App.css";
import Board from "./components/board";
import Keyboard from "./components/keyboard";
import CommandDisplay from "./components/commandDisplay";
import { ServerResponse, UserCommand } from "./models/model";

class App extends React.Component<any, any> {
  ws: WebSocket | null = null;

  constructor(props: any) {
    super(props);
    this.state = {
      isInGame: false,
      username: "",
      gameState: null,
      currentUsers: [],
      commandMessages: [],
    };
  }

  componentDidMount() {
    this.ws = new WebSocket("ws:localhost:8080/play");
    this.ws.onmessage = this.onWebSocketMessage.bind(this);
  }

  getCommandMessage(data: ServerResponse): string[] {
    const messages: string[] = [];
    const userId = data.UserCommandData.Metadata.UserId;
    switch (data.UserCommandData.CommandType) {
      case "JoinSession":
        messages.push(`${userId} joined the game.`);
        break;
      case "CreateSession":
        messages.push(`${userId} created a new game.`);
        break;
      case "CommitGuess":
        messages.push(`${userId} attempted a guess.`);
        break;
      case "BackspaceCharacter":
        messages.push(`${userId} deleted a character.`);
        break;
      case "AddCharacter":
        messages.push(
          `${userId} added character ${data.UserCommandData.AddCharacterCommandPayloadData?.Character}.`
        );
        break;
      case "ResetSession":
        messages.push(`${userId} reset the game.`);
        break;
    }
    if (data.ErrorMessage) {
      messages.push(`Error: ${data.ErrorMessage}`);
    }
    return messages;
  }

  onWebSocketMessage(event: MessageEvent) {
    const data: ServerResponse = JSON.parse(event.data) as ServerResponse;
    console.log(data);
    if (
      (data.UserCommandData.CommandType === "JoinSession" ||
        data.UserCommandData.CommandType === "CreateSession" ||
        data.UserCommandData.CommandType === "ResetSession") &&
      data.ErrorMessage === ""
    ) {
      this.setState({ ...this.state, isInGame: true });
    }
    this.setState({
      ...this.state,
      gameState: data.GameStateData,
      currentUsers: Object.keys(data.GameStateData.CurrentUsers),
      commandMessages: [
        ...this.state.commandMessages,
        ...this.getCommandMessage(data),
      ],
    });
  }

  onKeyboardClicked(cell: string) {
    if (cell === "ENTER") {
      const commitGuessRequest: UserCommand = {
        CommandType: "CommitGuess",
        Metadata: {
          SessionId: this.state.gameState.SessionId,
          UserId: this.state.username,
        },
      };
      this.ws?.send(JSON.stringify(commitGuessRequest));
      return;
    }
    if (cell === "DELETE") {
      const backspaceCharacterRequest: UserCommand = {
        CommandType: "BackspaceCharacter",
        Metadata: {
          SessionId: this.state.gameState.SessionId,
          UserId: this.state.username,
        },
      };
      this.ws?.send(JSON.stringify(backspaceCharacterRequest));
      return;
    }
    if (cell === "RESET GAME") {
      const resetSessionRequest: UserCommand = {
        CommandType: "ResetSession",
        Metadata: {
          SessionId: this.state.gameState.SessionId,
          UserId: this.state.username,
        },
      };
      this.ws?.send(JSON.stringify(resetSessionRequest));
      return;
    }
    const addCharacterRequest: UserCommand = {
      CommandType: "AddCharacter",
      Metadata: {
        SessionId: this.state.gameState.SessionId,
        UserId: this.state.username,
      },
      AddCharacterCommandPayloadData: {
        Character: cell,
      },
    };
    this.ws?.send(JSON.stringify(addCharacterRequest));
  }

  onNewGameButtonClicked() {
    if (!this.state.username) {
      alert("Username must not be empty");
      return;
    }
    const newGameRequest: UserCommand = {
      CommandType: "CreateSession",
      Metadata: {
        UserId: this.state.username,
      },
    };
    this.ws?.send(JSON.stringify(newGameRequest));
  }

  onJoinGameButtonClicked() {
    if (!this.state.username) {
      alert("Username must not be empty");
      return;
    }
    const sessionId = Number(window.prompt("Enter session ID:", ""));
    if (!sessionId) {
      window.alert("Session ID must be numbers");
      return;
    }
    const joinSessionRequest: UserCommand = {
      CommandType: "JoinSession",
      Metadata: {
        SessionId: sessionId,
        UserId: this.state.username,
      },
    };
    this.ws?.send(JSON.stringify(joinSessionRequest));
  }

  onUsernameChanged(event: ChangeEvent<HTMLInputElement>) {
    event.preventDefault();
    this.setState({ ...this.state, username: event.target.value });
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <div>CoWordle</div>
          {this.state.isInGame ? (
            <div>
              {this.state.gameState?.SessionId ? (
                <div className="session-id-container">
                  {" "}
                  Session ID: {this.state.gameState?.SessionId}
                </div>
              ) : (
                <></>
              )}
              <div className="app-board">
                <div className="app-board-left-container">
                  <div className="board-display-header">Board</div>
                  <Board gameState={this.state.gameState}></Board>
                </div>
                <div className="app-board-right-container">
                  <div className="command-display-header">Command Logs</div>
                  <CommandDisplay
                    commandMessages={this.state.commandMessages}
                  ></CommandDisplay>
                </div>
              </div>
              {this.state.gameState?.GameOver ? (
                <div className="game-over-message-container"> Game over. The word is {this.state.gameState?.Word}</div>
              ) : (
                <></>
              )}
              <div className="app-keyboard">
                <Keyboard
                  onKeyboardClicked={this.onKeyboardClicked.bind(this)}
                ></Keyboard>
              </div>
            </div>
          ) : (
            <div>
              <div>
                <input
                  type="text"
                  name="name"
                  className="username-input-container"
                  placeholder="Enter username"
                  value={this.state.username}
                  onChange={this.onUsernameChanged.bind(this)}
                />
              </div>
              <div>
                <span className="game-button" onClick={this.onNewGameButtonClicked.bind(this)}>
                  New game
                </span>
                <span className="game-button" onClick={this.onJoinGameButtonClicked.bind(this)}>
                  Join a game
                </span>
              </div>
            </div>
          )}
        </header>
      </div>
    );
  }
}

export default App;

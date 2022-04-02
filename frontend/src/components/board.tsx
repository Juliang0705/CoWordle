import React from "react";
import {GameState, Cell} from "../models/model";
import "./board.css";

class Board extends React.Component<any, any> {
  constructor(props: any) {
    super(props);
  }

  getCellDisplayValue(gameState: GameState, rowIndex: number, colIndex: number) {
    if (gameState.Grid[rowIndex][colIndex].Value) {
      return gameState.Grid[rowIndex][colIndex].Value;
    }
    const prevRowIsNotEmpty = rowIndex === 0 || gameState.Grid[rowIndex-1][0].Value !== '';
    if (gameState.Grid[rowIndex][0].Value === '' && prevRowIsNotEmpty) {
      if (colIndex < gameState.CurrentCharacters.length) {
        return gameState.CurrentCharacters[colIndex];
      }
    }
    return '';
  }

  render() {
    return (
      <div className="board">
        {this.props.gameState?.Grid.map((row: Array<Cell>, rowIndex: number) => {
          return (
            <div key={rowIndex} className="board-row">
              {row.map((cell: Cell, colIndex: number) => {
                return (
                  <div key={colIndex} className={["board-row-cell", cell.GuessType].join(' ')}>
                    <div>{this.getCellDisplayValue(this.props.gameState, rowIndex, colIndex)}</div>
                  </div>
                );
              })}
            </div>
          );
        })}
      </div>
    );
  }
}
export default Board;

import React from "react";
import "./keyboard.css";

class Keyboard extends React.Component<any, any> {
  keyboardRows: Array<Array<string>> = [];
  constructor(props: any) {
    super(props);
    this.keyboardRows = [
      ["Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"],
      ["A", "S", "D", "F", "G", "H", "J", "K", "L"],
      ["ENTER", "Z", "X", "C", "V", "B", "N", "M", "DELETE", "RESET GAME"],
    ];
  }
  render() {
    return (
      <div className="keyboard">
        {this.keyboardRows.map((row: Array<string>, index: number) => {
          return (
            <div key={index} className="keyboard-row">
              {row.map((cell: string, index: number) => {
                return (
                  <div
                    key={index}
                    className={[
                      "keyboard-row-cell",
                      cell === "ENTER" || cell === "DELETE" || cell === "RESET GAME"
                        ? "keyboard-row-cell-special"
                        : "",
                    ].join(" ")}
                    onClick={() => {this.props?.onKeyboardClicked(cell);}}
                  >
                    {cell}
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
export default Keyboard;

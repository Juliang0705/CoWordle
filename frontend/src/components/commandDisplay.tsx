import React from "react";
import "./commandDisplay.css";

class CommandDisplay extends React.Component<any, any> {
  constructor(props: any) {
    super(props);
  }
  render() {
    return (
      <div>
        <div className="command-display-board">
          {this.props.commandMessages
            ?.slice()
            .reverse()
            .map((message: string, index: number) => {
              return (
                <div
                  key={index}
                  className={[
                    "one-command",
                    message.startsWith("Error: ") ? "error-command" : "",
                  ].join(" ")}
                >
                  {message}
                </div>
              );
            })}
        </div>
      </div>
    );
  }
}

export default CommandDisplay;

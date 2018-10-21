import React from 'react';
import { Link } from 'react-router-dom';
const CreateSessionBtn = ({ ...props }) => (
  <span {...props}>
    <Link className="new-session" to="/Poker">
      Create Session
    </Link>
  </span>
);

export default class extends React.Component {
  render() {
    return (
      <span id="Home">
        <h1>Welcome!</h1>
        <CreateSessionBtn onClick={this.createSession} />
      </span>
    );
  }
}

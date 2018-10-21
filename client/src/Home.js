import React from 'react';
import { Link } from 'react-router-dom'
const CreateSessionBtn = ({ ...props }) => (
  <Link className="new-session" to="/Poker" {...props}>
    Create Session
  </Link>
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

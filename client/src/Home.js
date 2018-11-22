import React from 'react';
import { createNewSession } from './utils';

export default class extends React.Component {
  async onNewSessionClick() {
    const { ID, adminID } = await createNewSession();
    window.location.href = `/session?id=${ID}&adminId=${adminID}`;
  }
  render() {
    return (
      <span id="Home">
        <h1>Welcome!!</h1>
        <button onClick={this.onNewSessionClick}>Create Session</button>
      </span>
    );
  }
}

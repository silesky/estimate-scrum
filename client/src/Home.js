import React from 'react';
import { createNewSession } from './utils';

export default class extends React.Component {
  async onNewSessionClick() {
    const data = await createNewSession();
    const { ID, adminID } = data;
    window.location.href =
      `/session` +
      `?id=${ID}` +
      `&adminID=${adminID}`
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

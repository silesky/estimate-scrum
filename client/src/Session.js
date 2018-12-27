import React, { Component } from 'react';
import { getSession } from './utils';
import './App.css';

const mapEstimations = estimations =>
  Object.keys(estimations).map(u => ({
    username: u,
    estimate: estimations[u],
  }));

const createWebSocketConnection = (onMessageCb, { id, adminID }) => {
  const WS_URL = `ws://localhost:3333/ws?id=${id}&adminID=${adminID}`;
  const socket = new WebSocket(WS_URL);
  socket.addEventListener('open', function() {
    console.log('Websocket opened @ ', WS_URL);
  });
  socket.addEventListener('message', function(event) {
    const data = JSON.parse(event.data);
    console.log('Message from server ', data);
    onMessageCb(data);
  });
  return socket;
};

const Estimate = ({ name, estimate }) => <h4>{`${name}: ${estimate}`}</h4>;

const createMessage = (
  username,
  estimationValue,
  sessionID = 'abc123',
  issueID,
) => {
  // got some unexpected
  return JSON.stringify({ username, estimationValue, sessionID, issueID });
};

const CopyBox = ({ link }) => (
  <span id="CopyBox">
    <div className="copyBox">{link}</div>
  </span>
);

const AdminControlPanel = ({ isAdmin, setIssueTitle }) => {
  if (!isAdmin) return null;
  return (
    <div id="AdminControlPanel">
      <h2>ADMIN IS AUTHORIZED</h2>
      <label htmlFor="issueTitle">New Issue Title</label>
      <input
        type="text"
        id="issueTitle"
        onChange={e => setIssueTitle(e.target.value)}
      />
    </div>
  );
};

export default class extends Component {
  state = {
    currentUser: '',
    currentEstimate: null,
    estimations: {},
    isAdmin: false,
    error: false,
    // admin-only for setting
    issueTitle: '',
    issueID: '',
  };

  wsSubscription = data => {
    // callback
    console.assert(data.issueID, 'no issue id found!');
    console.assert(data.username, 'no username found!');
    this.setState({
      username: data.username,
      sessionID: data.sessionID,
      issueID: data.issueID,
      estimations: {
        ...this.state.estimations,
        [`${data.username}`]: data.estimationValue,
      },
    });
  };

  socket = createWebSocketConnection(this.wsSubscription, this.getParams());

  getParams() {
    const qp = new URLSearchParams(this.props.location.search);
    return {
      id: qp.get('id'),
      adminID: qp.get('adminID'),
      issueID: qp.get('issueID'),
    };
  }

  getNonAdminSessionLink() {
    const { id: sessionID } = this.getParams();
    return `${window.location.host}/session?id=${sessionID}`;
  }

  submitEstimation = () => {
    const { id: sessionID, issueID } = this.getParams();
    if (!sessionID) {
      console.error('no session ID found. should be ?id=1234');
      return;
    }
    const newEstimation = createMessage(
      this.state.currentUser,
      this.state.currentEstimate,
      sessionID,
      issueID,
    );
    this.socket.send(newEstimation);
  };

  setUser = currentUser => this.setState({ currentUser });
  setEstimate = currentEstimate => {
    const t = typeof currentEstimate;
    const est =
      t === 'string' || t === 'number' ? parseInt(currentEstimate, 10) : 0;
    this.setState({ currentEstimate: est });
  };
  setAdminStatus = isAdmin => this.setState({ isAdmin });
  setError = bool => this.setState({ error: bool });
  setIssueTitle = issueTitle => this.setState({ issueTitle });
  // http://localhost:3000/session?id=206f8d29-fa5a-4f0b-9051-6f7b4089967a
  async componentDidMount() {
    const { id, adminID } = this.getParams();
    // get any
    try {
      // TODO: wait till response comes back that confirms that user is admin or not
      const data = await getSession(id, adminID);
      console.log('data from /session', data);
      const { session, isAdmin } = data;
      this.setState({
        session,
        isAdmin,
      });
    } catch (err) {
      console.error(err, '__NO SESSION FOUND__');
      this.setError(true);
    }
  }
  render() {
    console.log('state:', this.state, 'props:', this.props);
    if (this.state.error) {
      const { id } = this.getParams();
      return <h3>{`Oops! No scrum session found with id: ${id}`}</h3>;
    }

    return (
      <div className="App">
        <AdminControlPanel
          isAdmin={this.state.isAdmin}
          setIssueTitle={this.setIssueTitle}
        />
        <CopyBox link={this.getNonAdminSessionLink()} />
        <h1>Scrum Session!</h1>
        {this.state.issueTitle && <h2>{this.state.issueTitle}</h2>}
        <div>
          <label htmlFor="username">Username</label>
          <input
            type="string"
            onChange={event => this.setUser(event.target.value)}
            id="username"
          />
        </div>
        <div>
          <label htmlFor="estimation">Estimation</label>
          <input
            onChange={event => this.setEstimate(event.target.value)}
            type="number"
            id="estimation"
          />
          <button id="submit" onClick={this.submitEstimation}>
            Submit
          </button>
          {mapEstimations(this.state.estimations).map((el, ind) => (
            <Estimate
              name={el.username}
              estimate={el.estimate}
              key={`${el.name}${el.estimate}${ind}`}
            />
          ))}
        </div>
      </div>
    );
  }
}

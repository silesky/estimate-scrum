import React, { Component } from 'react';
import { getSession } from './utils';
import { updateSession } from './utils/fetch';
import { pathOr } from 'ramda';
import './App.css';

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

const Estimate = ({ username, estimate }) => (
  <h4>{`${username}: ${estimate}`}</h4>
);

const Issue = ({ issue, isSelected }) => {
  const mapEstimations = estimations => {
    return Object.keys(estimations).map(u => ({
      username: u,
      estimate: estimations[u],
    }));
  };

  return (
    <React.Fragment>
      <h4 style={{color: isSelected ? 'red' : 'black'}}>IssueID: {issue.issueID}</h4>
      {mapEstimations(issue.estimations).map(estimate => (
        <Estimate
          username={estimate.username}
          estimate={estimate.estimate}
          key={estimate.username}
        />
      ))}
      <hr />
    </React.Fragment>
  );
};

const Issues = ({ issues, selectedIssue }) => {
  return (
    <React.Fragment>
      {issues.map(issue => (
        <Issue
          isSelected={issue.issueID === selectedIssue}
          issue={issue}
          key={issue.issueID}
        />
      ))}
    </React.Fragment>
  );
};
const createUserMessageEstimation = (
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

const AdminControlPanel = ({ isAdmin, setIssueTitle, setSelectedIssue }) => {
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
      <label htmlFor="selectedIssue">Selected Issue</label>
      <input
        type="text"
        id="selectedIssue"
        onChange={e => setSelectedIssue(e.target.value)}
      />
    </div>
  );
};

// {
//   "session": {
//     "dateCreated": "2018-12-31 06:23:47.193119 +0000 UTC",
//     "ID": "b8b1b9a2-1bb7-4b7f-8ebb-276e0c7e2aa9",
//     "storyPoints": [
//       1,
//       2,
//       3
//     ],
//     "issues": [
//       {
//         "issueTitle": "",
//         "issueID": "db1f7fd4-aff3-4495-80fa-ef842c7eda71",
//         "estimations": {
//           "13": 0,
//           "123": 1231,
//           "": 123,
//           "bar": 456,
//           "foo": 123
//         }
//       }
//     ],
//     "selectedIssue": "db1f7fd4-aff3-4495-80fa-ef842c7eda71"
//   },
//   "isAdmin": true
// }

const getIssues = pathOr([], ['session', 'issues']);
const getSelectedIssue = pathOr('', ['session', 'selectedIssue']);
const getIsAdmin = pathOr(false, ['session', 'isAdmin']);

export default class extends Component {
  state = {
    session: {
      dateCreated: '',
      ID: '',
      storyPoints: [],
      issues: [],
      selectedIssue: '',
    },
    issues: [],
    currentUser: '',
    currentEstimate: null,
    isAdmin: false,
    error: false,
    // admin-only for setting
    issueTitle: '',
    selectedIssue: '',
  };

  wsSubscription = data => {
    // callback

    console.assert(data.session.ID, 'no sessionID found!');
    this.setState({
      username: data.username,
      sessionID: data.session.ID,
      session: data.session,
      isAdmin: data.isAdmin,
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
    const newEstimation = createUserMessageEstimation(
      this.state.currentUser,
      this.state.currentEstimate,
      sessionID,
      issueID,
    );
    this.socket.send(newEstimation);
  };
  submitIssueTitle = title => {
    updateSession({ ...this.state.session, issueTitle: title });
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
  setSelectedIssue = issueID => this.setState({ selectedIssue: issueID });
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
        <pre>{JSON.stringify(this.state, null, 2)}</pre>
        <AdminControlPanel
          isAdmin={this.state.isAdmin}
          setIssueTitle={this.setIssueTitle}
          setSelectedIssue={this.setSelectedIssue}
        />
        <CopyBox link={this.getNonAdminSessionLink()} />
        <h1>Scrum Session!</h1>
        <h2>Selected issue: {getSelectedIssue(this.state)} </h2>
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
          <Issues
            selectedIssue={getSelectedIssue(this.state)}
            issues={getIssues(this.state)}
          />
        </div>
      </div>
    );
  }
}

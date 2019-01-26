import React, { Component } from 'react';
import { getSession, addEstimation, updateSession } from '../../utils';
import { pathOr } from 'ramda';
import { AdminPanel, StoryPointsSelector, Issue } from '../../Components';

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
const createEstimation = (
  username,
  estimationValue,
  sessionID = 'abc123',
  issueID,
) => ({ username, estimationValue, sessionID, issueID });

const CopyBox = ({ link }) => (
  <span id="CopyBox">
    <div className="copyBox">{link}</div>
  </span>
);

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
    currentUser: '',
    currentEstimate: null,
    isAdmin: false,
    error: false,
    // admin-only for setting
  };

  wsSubscription = data => {
    // callback
    console.assert(data.session.ID, 'no sessionID received from websocket!');
    this.setState({
      username: data.username,
      sessionID: data.session.ID,
      session: data.session,
    });
  };

  socket = createWebSocketConnection(this.wsSubscription, this.getParams());

  getParams() {
    const qp = new URLSearchParams(this.props.location.search);
    return {
      id: qp.get('id'),
      adminID: qp.get('adminID'),
    };
  }

  getNonAdminSessionLink() {
    const { id: sessionID } = this.getParams();
    return `${window.location.host}/session?id=${sessionID}`;
  }

  submitEstimation = async () => {
    const { id: sessionID } = this.getParams();
    const issueID = this.state.session.selectedIssue;
    const newEstimation = createEstimation(
      this.state.currentUser,
      this.state.currentEstimate,
      sessionID,
      issueID,
    );

    try {
      const res = await addEstimation(newEstimation);
      console.log(res);
    } catch (err) {
      console.warn(err);
    }
  };

  submitSessionUpdate = async () => {
    const { id, adminID } = this.getParams();
    try {
      await updateSession(id, adminID, this.state.session)
    } catch (err) {
      console.error(err);
    }
  };

  setIssueTitle = issueTitle => {
    const issues = this.state.session.issues.map(el => {
      if (el.issueID === this.state.session.selectedIssue) {
        return { ...el, issueTitle };
      }
      return el;
    });
    this.setState({
      session: {
        ...this.state.session,
        issues,
      },
    });
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
        <h1>Scrum Session!!</h1>
        <pre>{JSON.stringify(this.state, null, 2)}</pre>
        <AdminPanel
          isAdmin={this.state.isAdmin}
          setIssueTitle={this.setIssueTitle}
          setSelectedIssue={this.setSelectedIssue}
          submitSessionUpdate={this.submitSessionUpdate}
        />
        <CopyBox link={this.getNonAdminSessionLink()} />
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
        <StoryPointsSelector
          onChange={event => this.setEstimate(event.target.value)}
        />
        <button id="submit" onClick={this.submitEstimation}>
          Submit
        </button>
        <Issues
          selectedIssue={getSelectedIssue(this.state)}
          issues={getIssues(this.state)}
        />
      </div>
    );
  }
}

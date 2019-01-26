import React from 'react';
import styled from 'styled-components';
/* AdminControlPanel is only visible if the user is an admin. It allows the user to do the following:
- change current issue
- change current issue title (meh)
- change story points
*/
const AdminControlPanelContainer = styled.div`
  border: 1px solid black;
`;

const AdminControlPanel = ({ isAdmin, submitSessionUpdate, setIssueTitle, setSelectedIssue }) => {
  if (!isAdmin) return null;
  return (
    <AdminControlPanelContainer>
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
      <button id="new-issue">New Issue</button>
      <button onClick={() => submitSessionUpdate()}>SUBMIT</button>
    </AdminControlPanelContainer>
  );
};

export default AdminControlPanel;

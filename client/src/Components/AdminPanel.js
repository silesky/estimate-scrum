import React from 'react'
import StoryPointSelector from './StoryPointsSelector';

/* AdminControlPanel is only visible if the user is an admin. It allows the user to do the following:
- change current issue
- change current issue title (meh)
- change story points
*/
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
      <StoryPointSelector onChange={console.log} />
    </div>
  );
};

export default AdminControlPanel

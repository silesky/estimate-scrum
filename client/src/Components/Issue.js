import React from 'react'



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
      <button id="delete">DELETE ISSUE</button>
      <h4 style={{ color: isSelected ? 'red' : 'black' }}>
        IssueID: {issue.issueID}
      </h4>
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
export default Issue;

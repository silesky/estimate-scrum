import React from 'react';
const points = [[0, 1, 2, 3], [0, 3, 3, 4], [1, 2, 3, 4], [1, 2, 3, 4]];

export default props => (
  <select {...props}>
    {points.map(v => (
      <option value={v}>{v.join(' ')}</option>
    ))}
  </select>
);

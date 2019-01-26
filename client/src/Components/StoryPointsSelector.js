import React from 'react';
import { getFibonaccis } from '../utils'
const points = ["?", 0, ...getFibonaccis(10)]
export default props => (
  <select {...props}>
    {points.map(v => (
      <option key={v} value={v}>{v}</option>
    ))}
  </select>
);

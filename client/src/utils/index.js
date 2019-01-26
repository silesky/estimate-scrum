import * as R from 'ramda';
export * from './fetch';
export const toBool = bool => bool === 'true' || bool === true;
export const fibonacci = num => {
  if (num <= 1) return 1;
  return fibonacci(num - 1) + fibonacci(num - 2);
};
export const getFibonaccis = num => {
  return R.pipe(
    R.range(1),
    R.map(fibonacci),
  )(num + 1);
};

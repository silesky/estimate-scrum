import { fibonacci, getFibonaccis } from './index';

describe('fibonacci', () => {
  it('fibonacci should work', () => {
    expect(fibonacci(0)).toEqual(1);
    expect(fibonacci(1)).toEqual(1);
    expect(fibonacci(3)).toEqual(3);
    expect(fibonacci(4)).toEqual(5);
    expect(fibonacci(5)).toEqual(8);
    expect(fibonacci(6)).toEqual(13);
  });
  it('getFibonaccis should work', () => {
    expect(getFibonaccis(5)).toEqual([1, 2, 3, 5, 8]);
    expect(getFibonaccis(6)).toEqual([1, 2, 3, 5, 8, 13]);
  });
});

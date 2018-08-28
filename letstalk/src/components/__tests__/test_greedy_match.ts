import { greedyMatch } from '../AutocompleteInput';

it('returns null for no match', () => {
  expect(greedyMatch('abc', 'def')).toBeNull();
});

it('returns null for empty inputs', () => {
  expect(greedyMatch('', '')).toBeNull();
});

it('matches on end of word', () => {
  expect(greedyMatch('abcd', 'bcde')).toEqual([1, 4]);
});

it('matches on end of query', () => {
  expect(greedyMatch('abcd', 'bc')).toEqual([1, 3]);
});

it('matches earliest prefix', () => {
  expect(greedyMatch('abcd', 'abce')).toEqual([0, 3]);
});

it('matches earliest prefix with different case', () => {
  expect(greedyMatch('aBcd', 'Abce')).toEqual([0, 3]);
});

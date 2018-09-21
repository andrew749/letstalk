import { password, PASSWORD_ERROR } from '../index';

it('happy password', () => {
  expect(password('aB3$eF7*')).toBeUndefined();
});

it('too short password', () => {
  expect(password('aB3$eF7')).toEqual(PASSWORD_ERROR);
});

it('password no lowercase', () => {
  expect(password('AB3$EF7*')).toEqual(PASSWORD_ERROR);
});

it('password no uppercase', () => {
  expect(password('ab3$ef7*')).toEqual(PASSWORD_ERROR);
});

it('password no number', () => {
  expect(password('aBC$eFG*')).toEqual(PASSWORD_ERROR);
});

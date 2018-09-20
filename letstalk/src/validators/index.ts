const PASSWORD_ERROR = 'Password must contain an uppercase and lowercase letter, a number, and ' +
  'must be at least 8 characters long';

const required = (value: any) => (value ? undefined : 'Required')
const phoneNumber = (value: string) =>
  value && !/^(0|[1-9][0-9]{9})$/i.test(value)
    ? 'Invalid phone number, must be 10 digits'
    : undefined
const email = (value: string) =>
  value && !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(value)
    ? 'Invalid email address'
    : undefined;
const uwEmail = (value: string) =>
  value && !/^[A-Z0-9._%+-]+@(edu\.)?uwaterloo\.ca$/i.test(value)
    ? 'Invalid UW email address'
    : undefined;
const password = (value: string) => {
  const re = /(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}/;
  const elems = re.exec(value)
  return elems === null || elems.length === 0 ? PASSWORD_ERROR : undefined;
}

export {
  email,
  uwEmail,
  phoneNumber,
  required,
  password,
  PASSWORD_ERROR,
};

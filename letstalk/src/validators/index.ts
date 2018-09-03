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

export {
  email,
  uwEmail,
  phoneNumber,
  required,
};

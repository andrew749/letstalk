import Expo from 'expo';

// TODO: How to store this in some type of secrets file?
const API_KEY = '153025458705124';
const PERMISSIONS = ['public_profile', 'email', 'user_birthday'];

type AccessToken = string;

const fbLogin = async (): Promise<AccessToken> => {
  const { type, token } = await Expo.Facebook.logInWithReadPermissionsAsync(API_KEY, {
    permissions: PERMISSIONS,
  });

  if (type === 'success') return token;
  else return null;
}

export { fbLogin };

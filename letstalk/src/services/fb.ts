import Expo from 'expo';

// TODO: How to store this in some type of secrets file?
const API_KEY = '404473003315328';
const PERMISSIONS = ['public_profile', 'email', 'user_birthday'];

type AccessToken = {token: string, expires: number};

const fbLogin = async (): Promise<AccessToken> => {
  const { type, token, expires } = await Expo.Facebook.logInWithReadPermissionsAsync(API_KEY, {
    permissions: PERMISSIONS,
  });

  if (type === 'success') return { token, expires };
  else return null;
}

export { fbLogin };

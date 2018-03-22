// TODO: Figure out how to make this library work, since we would prefer to encrypt our token
// import * as Keychain from 'react-native-keychain';
import { AsyncStorage } from 'react-native';

import requestor from './requests';
import { SessionService, SessionToken, RemoteSessionService } from './session-service';

export class Auth {
  private sessionService: SessionService
  private sessionToken: SessionToken

  private static readonly tokenLocation = '@app/session';

  constructor(sessionService: SessionService) {
    this.sessionService = sessionService;
    this.sessionToken = null;
  }

  async getSessionToken(): Promise<SessionToken> {
    if (this.sessionToken === null) {
      this.sessionToken = await AsyncStorage.getItem(Auth.tokenLocation);
    }
    return this.sessionToken;
  }

  private async setSessionToken(sessionToken: SessionToken): Promise<void> {
    await AsyncStorage.setItem(Auth.tokenLocation, sessionToken);
    this.sessionToken = sessionToken;
  }

  async login(username: string, password: string, notificationToken?: string): Promise<void> {
    const sessionToken = await this.sessionService.login(
        username,
        password,
        notificationToken,
    );
    await this.setSessionToken(sessionToken);
  }

  async loginWithFb(notificationToken?: string): Promise<void> {
    const sessionToken = await this.sessionService.loginWithFb(notificationToken);
    console.log(sessionToken);
    if (sessionToken === null) return;
    await this.setSessionToken(sessionToken);
  }

  async logout(): Promise<void> {
    const sessionToken = await this.getSessionToken();
    AsyncStorage.removeItem(Auth.tokenLocation);
    await this.sessionService.logout(sessionToken);
  }
};

const auth = new Auth(new RemoteSessionService(requestor));

export default auth;

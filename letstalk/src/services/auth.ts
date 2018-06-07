// TODO: Figure out how to make this library work, since we would prefer to encrypt our token
// import * as Keychain from 'react-native-keychain';
import { AsyncStorage } from 'react-native';

import requestor from './requests';
import { SessionService, SessionToken, RemoteSessionService } from './session-service';
import {FORGOT_PASSWORD_ROUTE} from './constants';

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

  async login(email: string, password: string, notificationToken?: string): Promise<void> {
    const sessionToken = await this.sessionService.login(
        email,
        password,
        notificationToken,
    );
    await this.setSessionToken(sessionToken);
  }

  // Returns whether the login completed.
  async loginWithFb(notificationToken?: string): Promise<boolean> {
    const sessionToken = await this.sessionService.loginWithFb(notificationToken);
    if (sessionToken === null) return false;
    await this.setSessionToken(sessionToken);
    return true;
  }

  async logout(): Promise<void> {
    const sessionToken = await this.getSessionToken();
    AsyncStorage.removeItem(Auth.tokenLocation);
    await this.sessionService.logout(sessionToken);
  }

  async forgotPassword(email: string): Promise<void> {
    const resp = await requestor.post(FORGOT_PASSWORD_ROUTE, {"email": email});
  }

  async linkFB(): Promise<boolean> {
    const sessionToken = await this.getSessionToken();
    return await this.sessionService.linkFb(sessionToken);
  }
};

const auth = new Auth(new RemoteSessionService(requestor));

export default auth;

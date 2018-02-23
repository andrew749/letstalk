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

  async login(username: string, password: string): Promise<SessionToken> {
    const sessionToken = await this.sessionService.login(username, password);
    await this.setSessionToken(sessionToken);
    return this.getSessionToken();
  }

  async logout(): Promise<void> {
    const sessionToken = await this.getSessionToken();
    await this.sessionService.logout(sessionToken);
    AsyncStorage.removeItem(Auth.tokenLocation);
  }
};

const auth = new Auth(new RemoteSessionService(requestor));

export default auth;

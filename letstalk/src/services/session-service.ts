import { Requestor } from './requests';
import { FB_LOGIN_ROUTE, LOGIN_ROUTE, LOGOUT_ROUTE } from './constants';
import { fbLogin } from './fb';

export type SessionToken = string;

export interface SessionService {
  login(username: string, password: string, notificationToken?: string): Promise<SessionToken>;
  loginWithFb(notificationToken?: string): Promise<SessionToken>;
  logout(sessionToken: SessionToken): Promise<void>;
}

export class InvalidCredentialsError extends Error {
  constructor() {
    super('Invalid username or password');
  }
}

export class MockSessionService implements SessionService {
  private static readonly token: string = 'some-session-token'

  private static getRandomInt(min: number, max: number) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min)) + min;
  }

  async login(username: string, password: string, notificationToken?: string): Promise<SessionToken> {
    await new Promise(resolve => setTimeout(() => resolve(),
      MockSessionService.getRandomInt(100, 500)));
    if (username !== 'foo' || password !== 'bar') throw new InvalidCredentialsError();
    return MockSessionService.token;
  }

  async loginWithFb(): Promise<SessionToken> {
    return this.login('foo', 'bar');
  }

  async logout(sessionToken: SessionToken): Promise<void> {
    if (sessionToken !== MockSessionService.token) throw new Error('Invalid session token');
    // no-op
  }
}

export class RemoteSessionService implements SessionService {
  private requestor: Requestor

  constructor(requestor: Requestor) {
    this.requestor = requestor;
  }

  async login(email: string, password: string, notificationToken?: string): Promise<SessionToken> {
    const response = await this.requestor.post(LOGIN_ROUTE,
      { email, password, notificationToken });
    return response.sessionId;
  }

  async loginWithFb(notificationToken?: string): Promise<SessionToken> {
    const res = await fbLogin();
    if (!res) return null;
    const { token, expires } = res;
    const response = await this.requestor.post(FB_LOGIN_ROUTE, {
      token: token,
      expiry: expires,
      notificationToken: notificationToken
    });
    return response.sessionId;
  }

  async logout(sessionToken: SessionToken): Promise<void> {
    await this.requestor.post(LOGOUT_ROUTE, {}, sessionToken);
  }
}

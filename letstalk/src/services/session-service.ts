import { Requestor } from './requests';
import { LOGIN_ROUTE, LOGOUT_ROUTE } from './constants';

export type SessionToken = string | null;

export interface SessionService {
  login(username: string, password: string): Promise<SessionToken>;
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

  async login(username: string, password: string): Promise<SessionToken> {
    await new Promise(resolve => setTimeout(() => resolve(),
      MockSessionService.getRandomInt(100, 500)));
    if (username !== 'foo' || password !== 'bar') throw new InvalidCredentialsError();
    return MockSessionService.token;
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

  async login(username: string, password: string): Promise<SessionToken> {
    const response = await this.requestor.post(LOGIN_ROUTE,
      { userId: Number(username), password });
    return response.result.sessionId;
  }

  async logout(sessionToken: SessionToken): Promise<void> {
    // TODO: Actually do something here
    // await this.requestor.post(LOGOUT_ROUTE, { sessionToken });
  }
}


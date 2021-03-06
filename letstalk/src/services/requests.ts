import { BASE_URL } from './constants';
import { SessionToken } from './session-service';

export type ErrorTypes =
  | 'UNAUTHORIZED'
  | 'INVALID_REQUEST'
  | 'NO_NETWORK';

export interface APIError {
  readonly errorMsg: string;
  readonly errorType: ErrorTypes;
}

type Method = 'GET' | 'POST' | 'DELETE';

// TODO: this is incomplete, add stuff as you need. Or, try finding a good type def for fetch.
interface FetchOptions {
  method: Method;
  headers?: Headers;
  body?: string;
}

export class Requestor {

  private serverUrl: string;

  constructor(serverUrl: string) {
    this.serverUrl = serverUrl;
  }

  private async makeRequest(
    route: string,
    fetchParams: FetchOptions,
    sessionToken: SessionToken,
  ): Promise<any> {
    if (!fetchParams.headers) fetchParams.headers = new Headers();
    fetchParams.headers.append('sessionId', sessionToken);
    let response: any;
    try {
      response = await fetch(route, fetchParams);
    } catch(_e) {
      const e: APIError = { errorMsg: 'Could not connect to network', errorType: 'NO_NETWORK' };
      throw e;
    }
    if (!response.ok) return response.json().then((data: any) => {
      let errorType: ErrorTypes = 'INVALID_REQUEST';
      if (data.Error.code === 401) errorType = 'UNAUTHORIZED';
      const e: APIError = { errorMsg: data.Error.errorMsg || data.Error.message, errorType };
      throw e;
    });
    const data = await response.json();
    return data.Result;
  }

  async get(endpoint: string, sessionToken?: SessionToken): Promise<any> {
    const fetchParams: FetchOptions = {
      method: 'GET',
    };
    return this.makeRequest(this.serverUrl + endpoint, fetchParams, sessionToken);
  }

  private async withData(
    method: Method,
    endpoint: string,
    data: object,
    sessionToken?: SessionToken,
    headers?: Headers,
    isJson: boolean = false,
  ): Promise<any> {
    if (isJson && !headers) {
      const headers = new Headers({
        'Content-Type': 'application/json',
      });
    } else if (isJson && headers)  {
        headers.set('Content-Type', 'application/json');
    }

    const fetchParams: FetchOptions = {
      method,
      body: JSON.stringify(data),
      headers: headers,
    };
    return this.makeRequest(this.serverUrl + endpoint, fetchParams, sessionToken);
  }

  async post(
    endpoint: string,
    data: object,
    sessionToken?: SessionToken,
  ): Promise<any> {
    return this.withData('POST', endpoint, data, sessionToken, undefined, true);
  }

  async postFormData(
    endpoint: string,
    data: FormData,
    sessionToken?: SessionToken,
  ) {
    const headers = new Headers();
    headers.append('sessionId', sessionToken);
    return await fetch(this.serverUrl + endpoint, {
      method:'POST',
      body: data,
      headers: headers,
    });
  }

  async delete(endpoint: string, data: object, sessionToken?: SessionToken): Promise<any> {
    return this.withData('DELETE', endpoint, data, sessionToken);
  }
};

// singleton for requests
const requestor = new Requestor(BASE_URL);

export default requestor;

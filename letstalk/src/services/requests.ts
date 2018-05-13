import { BASE_URL } from './constants';
import { SessionToken } from './session-service';

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
    const response = await fetch(route, fetchParams);
    if (!response.ok) return response.json().then((data: any) => {
      throw new Error(data.Error.Message);
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

  async put(
      endpoint: string,
      data: object,
      sessionToken?: SessionToken,
  ): Promise<any> {
      return this.withData('PUT', endpoint, data, sessionToken, undefined, true);
  }

  async postBinary(
    endpoint: string,
    data: object,
    headers: Headers,
    sessionToken?: SessionToken,
  ) {
    return this.withData('POST', endpoint, data, sessionToken, headers, false)
  }

  async delete(endpoint: string, data: object, sessionToken?: SessionToken): Promise<any> {
    return this.withData('DELETE', endpoint, data, sessionToken);
  }
};

// singleton for requests
const requestor = new Requestor(BASE_URL);

export default requestor;

import { BASE_URL } from './constants';

export class Requestor {

  serverUrl: string;

  constructor(serverUrl: string) {
    this.serverUrl = serverUrl;
  }

  // TODO: Add types
  async _makeRequest(route: string, fetchParams: object): Promise<any> {
    const response = await fetch(route, fetchParams)
    if (!response.ok) return response.json().then((data: any) => {
      throw new Error(data.Error.Message);
    });
    return response.json();
  }

  async get(endpoint: string): Promise<any> {
    const fetchParams = {
      method: 'GET',
    };
    return this._makeRequest(this.serverUrl + endpoint, fetchParams);
  }

  async post(endpoint: string, data: object): Promise<any> {
    const fetchParams = {
      method: 'POST',
      body: JSON.stringify(data),
      headers: new Headers({
        'Content-Type': 'application/json',
      }),
    };
    return this._makeRequest(this.serverUrl + endpoint, fetchParams);
  }
};

// singleton for requests
const requestor = new Requestor(BASE_URL);

export default requestor;


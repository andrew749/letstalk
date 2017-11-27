// import fetch from 'cross-fetch';

class Requestor {

  serverUrl: string;

  constructor(serverUrl: string) {
    this.serverUrl = serverUrl;
  }

  _makeRequest(route: string, fetchParams: any) {
    return fetch(route, fetchParams).then((response: any) => {
      if (!response.ok) return response.json().then((data: any) => {
        throw new Error(data.error);
      });
      return response.json();
    }, (reason: any) => {
      throw reason;
    });
  }

  get(endpoint: string) {
    const fetchParams = {
      method: 'GET',
    };

    return this._makeRequest(this.serverUrl + endpoint, fetchParams)
  }

  post(endpoint: string, data: any[]) {
    const fetchParams = {
      method: 'POST',
      body: JSON.stringify(data),
      headers: new Headers({
        'Content-Type': 'application/json',
      }),
    };
    return this._makeRequest(this.serverUrl + endpoint, fetchParams)
  }
};

export default Requestor;

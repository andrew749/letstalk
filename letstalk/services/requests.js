import fetch from 'cross-fetch';

class Requestor {

  constructor(serverUrl) {
    this.serverUrl = serverUrl;
  }

  _makeRequest(route, fetchParams) {
    return fetch(route, fetchParams).then((response) => {
      if (!response.ok) return response.json().then((data) => {
        throw new Error(data.error);
      });
      return response.json();
    }, (reason) => {
      throw reason;
    });
  }

  get(endpoint) {
    const fetchParams = {
      method: 'GET',
    };

    return this._makeRequest(this.serverUrl + endpoint, fetchParams)
  }

  post(endpoint, data) {
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

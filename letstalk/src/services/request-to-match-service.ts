import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import {
  CredentialOptions,
  CredentialOrganization,
  CredentialPosition,
  ValidCredentialPair,
} from '../models/credential';
import { CREDENTIAL_OPTIONS_ROUTE } from './constants';

export interface RequestToMatchService {
  getCredentialOptions(): Promise<CredentialOptions>;
}

interface GetCredentialOptionsResponse {
  readonly validPairs: Array<ValidCredentialPair>;
  readonly organizations: Array<CredentialOrganization>;
  readonly positions: Array<CredentialPosition>;
}

export class RemoteRequestToMatchService implements RequestToMatchService {
  private requestor: Requestor;
  private auth: Auth;

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async getCredentialOptions(): Promise<CredentialOptions> {
    const response: GetCredentialOptionsResponse =
      await this.requestor.get(CREDENTIAL_OPTIONS_ROUTE);
    return {
      validPairs: Immutable.List(response.validPairs),
      organizations: Immutable.List(response.organizations),
      positions: Immutable.List(response.positions),
    };
  }
}

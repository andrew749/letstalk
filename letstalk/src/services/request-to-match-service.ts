import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import {
  CredentialWithId,
  CredentialPair,
  CredentialOptions,
  CredentialOrganization,
  CredentialPosition,
  CredentialFilterableElement,
  CredentialRequest,
  ValidCredentialPair,
} from '../models/credential';
import {
  CREDENTIAL_OPTIONS_ROUTE,
  CREDENTIAL_ROUTE,
  CREDENTIALS_ROUTE,
  CREDENTIAL_REQUEST_ROUTE,
  CREDENTIAL_REQUESTS_ROUTE,
} from './constants';

export interface RequestToMatchService {
  getCredentialOptions(): Promise<CredentialOptions>;
}

interface GetCredentialOptionsResponse {
  readonly validPairs: Array<ValidCredentialPair>;
  readonly organizations: Array<CredentialOrganization>;
  readonly positions: Array<CredentialPosition>;
}

type GetCredentialsResponse = Array<CredentialWithId>;

interface AddCredentialResponse {
  credentialId: number;
}

interface RemoveCredentialRequest {
  credentialId: number;
}

type GetCredentialRequestsResponse = Array<CredentialRequest>

interface AddCredentialRequestResponse {
  credentialRequestId: number;
}

interface RemoveCredentialRequestRequest {
  credentialRequestId: number;
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

  async getCredentials(): Promise<Immutable.List<CredentialWithId>> {
    const sessionToken = await auth.getSessionToken();
    const response: GetCredentialsResponse =
      await this.requestor.get(CREDENTIALS_ROUTE, sessionToken);
    return Immutable.List(response);
  }

  async addCredential(credential: CredentialPair): Promise<number> {
    const sessionToken = await auth.getSessionToken();
    const response: AddCredentialResponse =
      await this.requestor.post(CREDENTIAL_ROUTE, credential, sessionToken);
    return response.credentialId;
  }

  async removeCredential(credentialId: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const request: RemoveCredentialRequest = { credentialId };
    await this.requestor.delete(CREDENTIAL_ROUTE, request, sessionToken);
  }

  async getCredentialRequests(): Promise<Immutable.List<CredentialRequest>> {
    const sessionToken = await auth.getSessionToken();
    const response: GetCredentialRequestsResponse =
      await this.requestor.get(CREDENTIAL_REQUESTS_ROUTE, sessionToken);
    return Immutable.List(response);
  }

  async addCredentialRequest(credential: CredentialPair): Promise<number> {
    const sessionToken = await auth.getSessionToken();
    const response: AddCredentialRequestResponse =
      await this.requestor.post(CREDENTIAL_REQUEST_ROUTE, credential, sessionToken);
    return response.credentialRequestId;
  }

  async removeCredentialRequest(credentialRequestId: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const request: RemoveCredentialRequestRequest = { credentialRequestId };
    await this.requestor.delete(CREDENTIAL_REQUEST_ROUTE, request, sessionToken);
  }
}

const requestToMatchService = new RemoteRequestToMatchService(requestor, auth);
export default requestToMatchService;

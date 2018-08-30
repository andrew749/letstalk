import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { Credential, Credentials } from '../models/credential';
import {
  ALL_CREDENTIAL_ROUTE,
  CREDENTIAL_ROUTE,
  CREDENTIALS_ROUTE,
  CREDENTIAL_REQUEST_ROUTE,
  CREDENTIAL_REQUESTS_ROUTE,
  REMOVE_RTM_MATCHES_ROUTE,
  USER_SIMPLE_TRAIT_ROUTE,
  USER_SIMPLE_TRAIT_BY_NAME_ROUTE,
  USER_POSITION,
} from './constants';

type GetAllCredentialsResponse = Array<Credential>;
type GetCredentialsResponse = Array<Credential>;
type GetCredentialRequestsResponse = Array<Credential>;

interface AddCredentialRequest { name: string }
interface AddCredentialResponse { credentialId: number }
interface AddCredentialRequestRequest { name: string }
interface AddCredentialRequestResponse { credentialId: number }

interface AddUserSimpleTraitByIdRequest { simpleTraitId: number }
interface AddUserSimpleTraitByNameRequest { name: string }
interface AddUserPositionRequest {
  roleId?: number;
  roleName?: string;
  organizationId?: number;
  organizationName?: string;
  startDate: string;
  endDate?: string;
}
interface RemoveUserPositionRequest { userPositionId: number }
interface RemoveUserSimpleTraitRequest { userSimpleTraitId: number }

interface RemoveCredentialRequest { credentialId: number }
interface RemoveCredentialRequestRequest { credentialId: number }

export class RemoteRequestToMatchService {
  private requestor: Requestor;
  private auth: Auth;

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async getAllCredentials(): Promise<Credentials> {
    const response: GetAllCredentialsResponse =
      await this.requestor.get(ALL_CREDENTIAL_ROUTE);
    return Immutable.List(response);
  }

  async getCredentials(): Promise<Credentials> {
    const sessionToken = await auth.getSessionToken();
    const response: GetCredentialsResponse =
      await this.requestor.get(CREDENTIALS_ROUTE, sessionToken);
    return Immutable.List(response);
  }

  async addCredential(name: string): Promise<number> {
    const sessionToken = await auth.getSessionToken();
    const request: AddCredentialRequest = { name };
    const response: AddCredentialResponse =
      await this.requestor.post(CREDENTIAL_ROUTE, request, sessionToken);
    return response.credentialId;
  }

  async removeCredential(credentialId: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const request: RemoveCredentialRequest = { credentialId };
    await this.requestor.delete(CREDENTIAL_ROUTE, request, sessionToken);
  }

  async getCredentialRequests(): Promise<Credentials> {
    const sessionToken = await auth.getSessionToken();
    const response: GetCredentialRequestsResponse =
      await this.requestor.get(CREDENTIAL_REQUESTS_ROUTE, sessionToken);
    return Immutable.List(response);
  }

  async addCredentialRequest(name: string): Promise<number> {
    const sessionToken = await auth.getSessionToken();
    const request: AddCredentialRequestRequest = { name };
    const response: AddCredentialRequestResponse =
      await this.requestor.post(CREDENTIAL_REQUEST_ROUTE, request, sessionToken);
    return response.credentialId;
  }

  async removeCredentialRequest(credentialId: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const request: RemoveCredentialRequestRequest = { credentialId };
    await this.requestor.delete(CREDENTIAL_REQUEST_ROUTE, request, sessionToken);
  }

  async removeRtmMatches(userId: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const url = REMOVE_RTM_MATCHES_ROUTE + '/' + userId;
    await this.requestor.delete(url, null, sessionToken);
  }

  async addUserSimpleTraitById(id: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const req: AddUserSimpleTraitByIdRequest = { simpleTraitId: id };
    await this.requestor.post(USER_SIMPLE_TRAIT_ROUTE, req, sessionToken);
  }

  async addUserSimpleTraitByName(name: string): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const req: AddUserSimpleTraitByNameRequest = { name };
    await this.requestor.post(USER_SIMPLE_TRAIT_BY_NAME_ROUTE, req, sessionToken);
  }

  async removeUserSimpleTrait(id: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const req: RemoveUserSimpleTraitRequest = { userSimpleTraitId: id };
    await this.requestor.delete(USER_SIMPLE_TRAIT_ROUTE, req, sessionToken);
  }

  async addUserPosition(req: AddUserPositionRequest): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    await this.requestor.post(USER_POSITION, req, sessionToken);
  }

  async removeUserPosition(id: number): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const req: RemoveUserPositionRequest = { userPositionId: id };
    await this.requestor.delete(USER_POSITION, req, sessionToken);
  }
}

const requestToMatchService = new RemoteRequestToMatchService(requestor, auth);
export default requestToMatchService;

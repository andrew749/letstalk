import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import { BootstrapData, Cohort, Relationship } from '../models/bootstrap';
import auth, { Auth } from './auth';
import { BOOTSTRAP_ROUTE, COHORT_ROUTE, SIGNUP_ROUTE } from './constants';

interface SignupRequest {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  gender: string;
  birthday: number; // unix time
  password: string;
}

interface UpdateCohortRequest {
  cohortId: number;
}

export interface BootstrapResponse {
  readonly relationships: Immutable.List<Relationship>;
  readonly state: 'account_created' | 'account_setup' | 'account_matched';
  readonly cohort: Cohort;
};

export interface ProfileService {
  signup(request: SignupRequest): Promise<number>;
}

export class RemoteProfileService implements ProfileService {
  private requestor: Requestor
  private auth: Auth

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  async signup(request: SignupRequest): Promise<number> {
    const response = await this.requestor.post(SIGNUP_ROUTE, request);
    return response.userId;
  }

  async updateCohort(request: UpdateCohortRequest): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    await this.requestor.post(COHORT_ROUTE, request, sessionToken);
  }

  async bootstrap(): Promise<BootstrapData> {
    const sessionToken = await auth.getSessionToken();
    const response: BootstrapResponse = await this.requestor.get(BOOTSTRAP_ROUTE, sessionToken);
    // TODO, change to more sane types when response becomes camelCase
    const {
      relationships,
      state,
      cohort,
    } = response;
    return {
      relationships: Immutable.List(response.relationships),
      state,
      cohort,
    };
  }
}

const profileService = new RemoteProfileService(requestor, auth);
export default profileService;

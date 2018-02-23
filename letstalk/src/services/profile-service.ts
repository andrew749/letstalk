import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { COHORT_ROUTE, SIGNUP_ROUTE } from './constants';

interface SignupRequest {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  gender: string;
  birthday: number; // unix time
  password: string;
}

interface updateCohortRequest {
  cohortId: number;
}

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
    return response.result.userId;
  }

  async updateCohort(request: updateCohortRequest): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    await this.requestor.post(COHORT_ROUTE, request, sessionToken);
  }
}

const profileService = new RemoteProfileService(requestor, auth);
export default profileService;

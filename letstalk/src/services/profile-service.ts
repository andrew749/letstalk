import requestor, { Requestor } from './requests';
import { SIGNUP_ROUTE } from './constants';

interface SignupRequest {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  gender: string;
  birthday: number; // unix time
  password: string;
}

export interface ProfileService {
  signup(request: SignupRequest): Promise<number>;
}

export class RemoteProfileService implements ProfileService {
  private requestor: Requestor

  constructor(requestor: Requestor) {
    this.requestor = requestor;
  }

  async signup(request: SignupRequest): Promise<number> {
    const response = await this.requestor.post(SIGNUP_ROUTE, request);
    return response.result.userId;
  }
}

const profileService = new RemoteProfileService(requestor);
export default profileService;

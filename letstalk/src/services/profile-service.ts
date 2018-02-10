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

  /*
  6   {                                                                                                                | 21
  5     'first_name': string,                                                                                          | 20 /**
  4     'last_name': string,                                                                                           | 19  * Method to get called in login route
  3     'email': string,                                                                                               | 18  * ```
  2                 'phone_number' string(optional),                                                                   | 17  *  {"user_id": string, "password": string}
  1                 'gender': string,                                                                                  | 16  * ```
28                  'birthday': date,                                                                                  | 15
  1     'password': string,                                                                                            | 14 func LoginUser(c *ctx.Context) errs.Error {
  2   }
  */

export interface ProfileService {
  signup(request: SignupRequest): Promise<number>;
}

export class RemoteProfileService implements ProfileService {
  private requestor: Requestor

  constructor(requestor: Requestor) {
    this.requestor = requestor;
  }

  async signup(request: SignupRequest): Promise<number> {
    const response = await this.requestor.post(SIGNUP_ROUTE,
      {
        ...request,
        first_name: request.firstName,
        last_name: request.lastName,
        phone_number: request.phoneNumber,
      });
    // TODO: typing
    return response.Result.UserId;
  }
}

const profileService = new RemoteProfileService(requestor);
export default profileService;

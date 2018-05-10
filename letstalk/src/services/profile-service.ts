import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import {
  BootstrapData,
  Cohort,
  OnboardingStatus,
  Relationship,
  UserData,
  UserState,
  UserType,
} from '../models';
import auth, { Auth } from './auth';
import {
  BOOTSTRAP_ROUTE,
  COHORT_ROUTE,
  SIGNUP_ROUTE,
  USER_VECTOR_ROUTE,
} from './constants';

export interface SignupRequest {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  gender: string;
  birthday: number; // unix time
  password: string;
  profilePic?: string;
}

interface UpdateCohortRequest {
  cohortId: number;
}

export interface PersonalityVector {
  sociable: number;
  hardworking: number;
  ambitious: number;
  energetic: number;
  carefree: number;
  confident: number;
}

export enum UserVectorPreferenceType {
  PREFERENCE_TYPE_ME = 0,
  PREFERENCE_TYPE_YOU
}

type UpdateVectorRequest = PersonalityVector & {
  readonly preferenceType: UserVectorPreferenceType;
};

export interface BootstrapResponse {
  readonly relationships: Array<Relationship>;
  readonly state: UserState;
  readonly cohort: Cohort;
  readonly me: UserData;
  readonly onboardingStatus: OnboardingStatus;
};

interface OnboardingUpdateResponse {
  readonly message: string;
  readonly onboardingStatus: OnboardingStatus;
}

export interface ProfileService {
  signup(request: SignupRequest): Promise<number>;
  updateCohort(request: UpdateCohortRequest): Promise<OnboardingStatus>;
  updateVector(
    preferenceType: UserVectorPreferenceType,
    vector: PersonalityVector
  ): Promise<OnboardingStatus>;
  bootstrap(): Promise<BootstrapData>;
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

  async updateCohort(request: UpdateCohortRequest): Promise<OnboardingStatus> {
    const sessionToken = await this.auth.getSessionToken();
    const response: OnboardingUpdateResponse = await this.requestor.post(
      COHORT_ROUTE, request, sessionToken);
    return response.onboardingStatus;
  }

  async updateVector(
    preferenceType: UserVectorPreferenceType,
    vector: PersonalityVector
  ): Promise<OnboardingStatus> {
    const sessionToken = await this.auth.getSessionToken();
    const request: UpdateVectorRequest = {
      ...vector,
      preferenceType,
    };
    const response: OnboardingUpdateResponse = await this.requestor.post(
      USER_VECTOR_ROUTE, request, sessionToken);
    return response.onboardingStatus;
  }

  async bootstrap(): Promise<BootstrapData> {
    const sessionToken = await this.auth.getSessionToken();
    const response: BootstrapResponse = await this.requestor.get(BOOTSTRAP_ROUTE, sessionToken);
    const {
      relationships,
      me,
    } = response;
    return {
      ...response,
      relationships: Immutable.List(relationships),
      me: {
        ...me,
        birthdate: new Date(me.birthdate),
      },
    };
  }


}

export const profileService = new RemoteProfileService(requestor, auth);
export default profileService;

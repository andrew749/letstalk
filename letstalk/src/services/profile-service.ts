import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import {
  BootstrapData,
  Cohort,
  OnboardingStatus,
  ProfileData,
  Relationship,
  UserState,
  UserType,
} from '../models';
import {
  UserPersonalInfo,
} from '../models/user';
import auth, { Auth } from './auth';
import {
  BOOTSTRAP_ROUTE,
  COHORT_ROUTE,
  ME_ROUTE,
  SIGNUP_ROUTE,
  USER_VECTOR_ROUTE,
  PROFILE_EDIT_ROUTE,
} from './constants';

export interface SignupRequest {
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
  readonly phoneNumber: string;
  readonly gender: number;
  readonly birthdate: number; // unix time
  readonly password: string;
  readonly profilePic?: string;
}

interface UpdateCohortRequest {
  readonly cohortId: number;
  readonly mentorshipPreference: number;
}

export interface PersonalityVector {
  readonly sociable: number;
  readonly hardworking: number;
  readonly ambitious: number;
  readonly energetic: number;
  readonly carefree: number;
  readonly confident: number;
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
  readonly me: UserPersonalInfo;
  readonly onboardingStatus: OnboardingStatus;
};

export interface ProfileEditRequest {
  readonly firstName: string;
  readonly lastName: string;
  readonly gender: number;
  readonly birthdate: number; // unix time
  readonly phoneNumber: string | null;
  readonly cohortId: number;
}

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

  async profileEdit(request: ProfileEditRequest): Promise<void> {
    const sessionToken = await this.auth.getSessionToken();
    await this.requestor.post(PROFILE_EDIT_ROUTE, request, sessionToken);
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

  async me(): Promise<ProfileData> {
    const sessionToken = await this.auth.getSessionToken();
    const response: ProfileData = await this.requestor.get(ME_ROUTE, sessionToken);
    return {
      ...response,
      birthdate: new Date(response.birthdate),
    };
  }

  static getProfilePicUrl(userId: string): string {
      return `https://s3.amazonaws.com/hive-user-profile-pictures/${userId}`;
  }

  static async getProfilePicForUser(userId: string): Promise<Blob> {
    const profilePicUrl = this.getProfilePicUrl(userId);
    // not a standard get
    const response = await fetch(profilePicUrl);
    console.log(response);
    const data = await response.body;
    const blob = (await data.getReader().read() as Blob);
    return blob;
  }

}

export const profileService = new RemoteProfileService(requestor, auth);
export default profileService;

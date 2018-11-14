import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import {
  BootstrapData,
  Cohort,
  OnboardingStatus,
  ProfileData,
  Relationship,
  UserState,
} from '../models';
import {
  MatchProfileData,
} from '../models/profile';
import {
  UserAdditionalData,
  UserPersonalInfo,
} from '../models/user';
import { Notification, NotifState } from '../models/notification';
import auth, { Auth } from './auth';
import {
  BOOTSTRAP_ROUTE,
  COHORT_ROUTE,
  COHORTS_ROUTE,
  MATCH_PROFILE_ROUTE,
  ME_ROUTE,
  NOTIFICATIONS_ROUTE,
  NOTIFICATION_ROUTE,
  NOTIFICATIONS_UPDATE_STATE_ROUTE,
  SIGNUP_ROUTE,
  USER_VECTOR_ROUTE,
  PROFILE_EDIT_ROUTE,
  PROFILE_PIC_ROUTE,
} from './constants';

export interface SignupRequest {
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
  readonly phoneNumber: string;
  readonly gender: number;
  readonly birthdate: string;
  readonly password: string;
  readonly profilePic?: string;
}

interface UpdateCohortRequest extends UserAdditionalData {
  readonly cohortId: number;
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

export interface ProfileEditRequest extends UserAdditionalData {
  readonly firstName: string;
  readonly lastName: string;
  readonly gender: number;
  readonly birthdate: string;
  readonly phoneNumber: string | null;
  readonly cohortId: number;
}

interface OnboardingUpdateResponse {
  readonly message: string;
  readonly onboardingStatus: OnboardingStatus;
}

interface NotificationRes {
  notificationId: number;
  title: string;
  userId: number;
  type: string;
  state: string;
  data: object;
  message: string;
  thumbnail: string;
  timestamp: string;
  link: string;
}

interface UpdateNotificationStateRequest {
  notificationIds: Array<number>;
  state: string;
}

export interface ProfileService {
  signup(request: SignupRequest): Promise<number>;
  updateCohort(request: UpdateCohortRequest): Promise<OnboardingStatus>;
  updateVector(
    preferenceType: UserVectorPreferenceType,
    vector: PersonalityVector
  ): Promise<OnboardingStatus>;
  bootstrap(): Promise<BootstrapData>;
  getProfilePicUrl(userId: string): Promise<string>;
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
    const response: BootstrapData = await this.requestor.get(BOOTSTRAP_ROUTE, sessionToken);

    return {
      ...response,
      connections: {
        outgoingRequests: Immutable.List(response.connections.outgoingRequests),
        incomingRequests: Immutable.List(response.connections.incomingRequests),
        mentors: Immutable.List(response.connections.mentors),
        mentees: Immutable.List(response.connections.mentees),
        peers: Immutable.List(response.connections.peers),
      },
    };
  }

  async getAllCohorts(): Promise<Immutable.List<Cohort>> {
    const response: Array<Cohort> = await this.requestor.get(COHORTS_ROUTE);
    return Immutable.List(response);
  }

  // NOTE: Could actually just return `ProfileData`, but this makes the types work
  private async getProfile(url: string): Promise<MatchProfileData> {
    const sessionToken = await this.auth.getSessionToken();
    const response: MatchProfileData = await this.requestor.get(url, sessionToken);
    return {
      ...response,
      userPositions: Immutable.List(response.userPositions),
      userSimpleTraits: Immutable.List(response.userSimpleTraits),
      userGroupSurveys: Immutable.List(response.userGroupSurveys).map(groupSurvey => {
        return {
          ...groupSurvey,
          survey: {
            ...groupSurvey.survey,
            responses: groupSurvey.survey.responses && Immutable.Map(groupSurvey.survey.responses),
            questions: Immutable.List(groupSurvey.survey.questions),
          },
        };
      }).toList(),
    };
  }

  async me(): Promise<ProfileData> {
    return this.getProfile(ME_ROUTE);
  }

  async matchProfile(userId: number): Promise<MatchProfileData> {
    const url = MATCH_PROFILE_ROUTE + '/' + userId;
    return this.getProfile(url);
  }

  async getProfilePicUrl(userId: string): Promise<string> {
    const sessionToken = await this.auth.getSessionToken();
    const profileRequest = `${PROFILE_PIC_ROUTE}?userId=${userId}`;
    const response: {profilePic: string} = await this.requestor.get(profileRequest, sessionToken);
    return response.profilePic;
  }

  async getNotifications(limit: number, past?: number): Promise<Immutable.List<Notification>> {
    const sessionToken = await this.auth.getSessionToken();
    let notificationsRequest = `${NOTIFICATIONS_ROUTE}?limit=${limit}`;
    if (!!past) {
      notificationsRequest = notificationsRequest + `&past=${past}`;
    }

    const response: Array<NotificationRes> =
      await this.requestor.get(notificationsRequest, sessionToken);

    return Immutable.List(response.map(notifRes => {
      return {
        ...notifRes,
        timestamp: new Date(notifRes.timestamp),
      } as Notification;
    }));
  }

  async getNotificationForId(notificationId: number): Promise<Notification> {
    const sessionToken = await this.auth.getSessionToken();
    let request = `${NOTIFICATION_ROUTE}/${notificationId}`;

    const response: NotificationRes =
      await this.requestor.get(request, sessionToken);
    return {
      ...response,
      timestamp: new Date(response.timestamp),
    } as Notification;
  }

  async updateNotificationState(
    notificationIds: Immutable.List<number>,
    state: NotifState,
  ): Promise<void> {
    const sessionToken = await this.auth.getSessionToken();
    const request: UpdateNotificationStateRequest = {
      notificationIds: notificationIds.toJS(),
      state,
    };
    await this.requestor.post(NOTIFICATIONS_UPDATE_STATE_ROUTE, request, sessionToken);
  }
}

export const profileService = new RemoteProfileService(requestor, auth);
export default profileService;

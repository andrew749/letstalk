export enum GenderId {
  Other = 0,
  Female,
  Male,
};

export function genderIdToString(genderId: GenderId): string {
  switch (genderId) {
    case GenderId.Other:
      return 'other';
    case GenderId.Male:
      return 'male';
    case GenderId.Female:
      return 'female';
    default:
      const _: never = genderId;
  }
}

export const USER_TYPE_MENTOR   = 1;
export const USER_TYPE_MENTEE   = 2;
export const USER_TYPE_ASKER    = 3;
export const USER_TYPE_ANSWERER = 4;
export const USER_TYPE_UNKNOWN  = -1;

export const MENTORSHIP_PREFERENCE_MENTOR = 1;
export const MENTORSHIP_PREFERENCE_MENTEE = 2;

export interface UserPersonalInfo {
  readonly userId: number
	readonly firstName: string
	readonly lastName: string
	readonly gender: GenderId
	readonly birthdate: Date
  readonly secret: string
}

export interface UserAdditionalData {
  readonly mentorshipPreference: number | null;
  readonly bio: string | null;
  readonly hometown: string | null;
}

export interface UserContactInfo {
	readonly email: string
  readonly phoneNumber: string | null;
  readonly fbId: string | null;
}

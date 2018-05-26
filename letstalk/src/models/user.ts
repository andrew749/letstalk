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

export const USER_TYPE_MENTOR  = 'user_type_mentor';
export const USER_TYPE_MENTEE  = 'user_type_mentee';
export const USER_TYPE_UNKNOWN = 'user_type_unknown';

export const MENTORSHIP_PREFERENCE_MENTOR = 1;
export const MENTORSHIP_PREFERENCE_MENTEE = 2;

export type UserType = 'user_type_mentor' | 'user_type_mentee' | 'user_type_unknown';

export interface UserPersonalInfo {
  readonly userId: number
	readonly firstName: string
	readonly lastName: string
	readonly gender: GenderId
	readonly birthdate: Date
}

export interface UserContactInfo {
	readonly email: string
  readonly phoneNumber: string | null;
  readonly fbId: string | null;
}

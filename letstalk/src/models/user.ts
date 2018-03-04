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

export interface UserData {
	readonly userId: number
	readonly firstName: string
	readonly lastName: string
	readonly email: string
	readonly gender: GenderId
	readonly birthdate: Date
}

import { Cohort } from './cohort';
import { UserAdditionalData, UserContactInfo, UserPersonalInfo } from './user';

export interface ProfileData extends UserAdditionalData, UserPersonalInfo, UserContactInfo,
  Cohort, UserContactInfo {
}

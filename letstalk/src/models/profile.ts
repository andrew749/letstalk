import { Cohort } from './cohort';
import { UserContactInfo, UserPersonalInfo } from './user';

export interface ProfileData extends UserPersonalInfo, UserContactInfo, Cohort {
}

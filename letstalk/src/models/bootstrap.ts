import Immutable from 'immutable';

import { Cohort } from './cohort';
import { UserPersonalInfo, UserType } from './user';
import { OnboardingStatus } from './onboarding';

export interface Relationship {
  readonly userId: number;
  readonly userType: UserType;
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
  readonly fbId: string | null;
  readonly phoneNumber: string | null;
}

export type UserState =
  | 'account_created'
  | 'account_setup'
  | 'account_matched';

export interface BootstrapData {
  readonly relationships: Immutable.List<Relationship>;
  readonly state: UserState;
  readonly cohort: Cohort;
  readonly onboardingStatus: OnboardingStatus;
};

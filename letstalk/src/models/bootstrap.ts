import Immutable from 'immutable';

import { Cohort } from './cohort';
import { UserData, UserType } from './user';
import { OnboardingStatus } from './onboarding';

export interface Relationship {
  readonly userId: number;
  readonly userType: UserType;
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
}

export type UserState =
  | 'account_created'
  | 'account_setup'
  | 'account_matched';

export interface BootstrapData {
  readonly relationships: Immutable.List<Relationship>;
  readonly state: UserState;
  readonly cohort: Cohort;
  readonly me: UserData;
  readonly onboardingStatus: OnboardingStatus;
};

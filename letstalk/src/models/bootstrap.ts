import Immutable from 'immutable';

import { Cohort } from './cohort';
import { UserPersonalInfo } from './user';
import { OnboardingStatus } from './onboarding';

export interface Relationship {
  readonly userId: number;
  readonly userType: number;
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
  readonly fbId: string | null;
  readonly fbLink: string | null;
  readonly phoneNumber: string | null;
  readonly cohort: Cohort | null;
  readonly description: string | null;
  readonly matchingState: MatchingState;
}

export enum MatchingState {
  Unknown = 0,
  Unverified,
  Verified,
  Expired,
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

import { UserType } from './user';

export const ONBOARDING_COHORT     = 'onboarding_cohort'     // get cohort info
export const ONBOARDING_VECTOR_ME  = 'onboarding_vector_me'  // get my personality vector
export const ONBOARDING_VECTOR_YOU = 'onboarding_vector_you' // get personality vector for others
export const ONBOARDING_DONE       = 'onboarding_done'       // finished

export type OnboardingState =
  | 'onboarding_cohort'
  | 'onboarding_vector_me'
  | 'onboarding_vector_you'
  | 'onboarding_done';

export interface OnboardingStatus {
  readonly state: OnboardingState;
  readonly userType: UserType;
}

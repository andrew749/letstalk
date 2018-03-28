import { Action } from 'redux'

import { OnboardingStatus } from '../../models';

export enum TypeKeys {
  SET_ONBOARDING_STATUS = 'ONBOARDING/SET_ONBOARDING_STATUS',
}

export interface SetOnboardingStatusAction extends Action {
  readonly type: TypeKeys.SET_ONBOARDING_STATUS;
  readonly onboardingStatus: OnboardingStatus;
}

export function setOnboardingStatusAction(
  onboardingStatus: OnboardingStatus
): SetOnboardingStatusAction {
  return {
    type: TypeKeys.SET_ONBOARDING_STATUS,
    onboardingStatus,
  };
}

export type ActionTypes =
  | SetOnboardingStatusAction

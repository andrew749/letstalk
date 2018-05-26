import Immutable from 'immutable';

import {
  ActionTypes,
  TypeKeys,
} from './actions';
import { ONBOARDING_DONE } from '../../models/onboarding';
import { USER_TYPE_UNKNOWN } from '../../models/user';
import {
  Cohort,
  OnboardingStatus,
} from '../../models';

export interface State {
  readonly onboardingStatus: OnboardingStatus;
}

// TODO: Fetch cohorts from server
const initialState: State = {
  onboardingStatus: {
    state: ONBOARDING_DONE,
    userType: USER_TYPE_UNKNOWN,
  },
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.SET_ONBOARDING_STATUS:
      return {
        ...state,
        onboardingStatus: action.onboardingStatus,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

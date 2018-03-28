import Immutable from 'immutable';

import {
  ActionTypes,
  TypeKeys,
} from './actions';
import {
  Cohort,
  ONBOARDING_DONE,
  OnboardingStatus,
  USER_TYPE_UNKNOWN,
} from '../../models';

export interface State {
  readonly cohorts: Immutable.List<Cohort>;
  readonly onboardingStatus: OnboardingStatus;
}

// TODO: Fetch cohorts from server
const initialState: State = {
  cohorts: Immutable.List([
    { cohortId: 1, programId: 'SOFTWARE_ENGINEERING', sequenceId: '8STREAM', gradYear: 2018},
    { cohortId: 2, programId: 'COMPUTER_ENGINEERING', sequenceId: '8STREAM', gradYear: 2018},
    { cohortId: 3, programId: 'COMPUTER_ENGINEERING', sequenceId: '4STREAM', gradYear: 2018},
  ]),
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
  }
};

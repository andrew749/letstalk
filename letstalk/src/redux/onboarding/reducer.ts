import Immutable from 'immutable';

import {
  Step,
  ActionTypes,
  TypeKeys,
} from './actions';
import {
  Cohort,
} from '../../models/cohort';

export interface State {
  cohorts: Immutable.List<Cohort>;
  step: Step;
}

// TODO: Fetch cohorts from server
const initialState: State = {
  cohorts: Immutable.List([
    { cohortId: 1, programId: 'SOFTWARE_ENGINEERING', sequenceId: '8STREAM', gradYear: 2018},
    { cohortId: 2, programId: 'COMPUTER_ENGINEERING', sequenceId: '8STREAM', gradYear: 2018},
    { cohortId: 3, programId: 'COMPUTER_ENGINEERING', sequenceId: '4STREAM', gradYear: 2018},
  ]),
  step: Step.COHORT,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.SET_STEP:
      const { step } = action;
      return {
        ...state,
        step,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

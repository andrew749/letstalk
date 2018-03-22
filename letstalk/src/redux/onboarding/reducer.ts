import {
  Step,
  ActionTypes,
  TypeKeys,
} from './actions';

export interface State {
  step: Step;
}

const initialState: State = {
  step: Step.MY_VECTOR,
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

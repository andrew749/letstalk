import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';

import { reducer as bootstrap, State as BootstrapState } from './bootstrap/reducer';
import { reducer as onboarding, State as OnboardingState } from './onboarding/reducer';

export interface RootState {
  bootstrap: BootstrapState;
  form: FormState;
  onboarding: OnboardingState;
};

const reducers = {
  bootstrap,
  form,
  onboarding,
};

const appReducer: Reducer<RootState> = combineReducers<RootState>(reducers);

export default appReducer;

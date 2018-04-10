import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';

import { reducer as bootstrap, State as BootstrapState } from './bootstrap/reducer';
import { reducer as onboarding, State as OnboardingState } from './onboarding/reducer';
import {
  reducer as credentialOptions,
  State as CredentialOptionsState,
} from './credential-options/reducer';

export interface RootState {
  bootstrap: BootstrapState;
  credentialOptions:CredentialOptionsState
  form: FormState;
  onboarding: OnboardingState;
};

const reducers = {
  bootstrap,
  credentialOptions,
  form,
  onboarding,
};

const appReducer: Reducer<RootState> = combineReducers<RootState>(reducers);

export default appReducer;

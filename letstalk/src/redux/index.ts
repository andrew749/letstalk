import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';
import { toastReducer as toast } from 'react-native-redux-toast';

import { reducer as bootstrap, State as BootstrapState } from './bootstrap/reducer';
import { reducer as onboarding, State as OnboardingState } from './onboarding/reducer';
import { reducer as credentials, State as CredentialsState } from './credentials/reducer';
import {
  reducer as credentialOptions,
  State as CredentialOptionsState,
} from './credential-options/reducer';

export interface RootState {
  bootstrap: BootstrapState;
  credentialOptions: CredentialOptionsState,
  credentials: CredentialsState,
  form: FormState;
  onboarding: OnboardingState;
};

const reducers = {
  bootstrap,
  credentialOptions,
  credentials,
  form,
  toast,
  onboarding,
};

const appReducer: Reducer<RootState> = combineReducers<RootState>(reducers);

export default appReducer;

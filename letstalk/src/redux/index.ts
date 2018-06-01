import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';
import { toastReducer as toast } from 'react-native-redux-toast';

import { reducer as bootstrap, State as BootstrapState } from './bootstrap/reducer';
import { reducer as onboarding, State as OnboardingState } from './onboarding/reducer';
import { reducer as credentials, State as CredentialsState } from './credentials/reducer';
import { reducer as matchProfile, State as MatchProfileState } from './match-profile/reducer';
import { reducer as profile, State as ProfileState } from './profile/reducer';
import {
  reducer as credentialRequests,
  State as CredentialRequestsState,
} from './credential-requests/reducer';
import {
  reducer as credentialOptions,
  State as CredentialOptionsState,
} from './credential-options/reducer';

export interface RootState {
  bootstrap: BootstrapState;
  credentialOptions: CredentialOptionsState,
  credentialRequests: CredentialRequestsState,
  credentials: CredentialsState,
  form: FormState;
  onboarding: OnboardingState;
  matchProfile: MatchProfileState,
  profile: ProfileState;
};

const reducers = {
  bootstrap,
  credentialOptions,
  credentialRequests,
  credentials,
  form,
  toast,
  onboarding,
  matchProfile,
  profile,
};

const appReducer: Reducer<RootState> = combineReducers<RootState>(reducers);

export default appReducer;

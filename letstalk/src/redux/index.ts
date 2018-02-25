import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';

import { reducer as bootstrap, State as BootstrapState } from './bootstrap/reducer';

export interface RootState {
  form: FormState;
  bootstrap: BootstrapState;
};

const appReducer: Reducer<RootState> = combineReducers<RootState>({ form, bootstrap });

export default appReducer;

import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';

export interface RootState {
  form: FormState;
};

const appReducer: Reducer<RootState> = combineReducers<RootState>({ form });

export default appReducer;

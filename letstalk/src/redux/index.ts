import { combineReducers, Reducer } from 'redux';
import { reducer as form, FormState } from 'redux-form';
import { toastReducer as toast } from 'react-native-redux-toast';

import { reducer as bootstrap, State as BootstrapState } from './bootstrap/reducer';
import { reducer as cohorts, State as CohortsState } from './cohorts/reducer';
import { reducer as matchProfile, State as MatchProfileState } from './match-profile/reducer';
import { reducer as notifications, State as NotificationsState } from './notifications/reducer';
import { reducer as profile, State as ProfileState } from './profile/reducer';
import { reducer as userSearch, State as UserSearchState } from './user-search/reducer';
import { reducer as searchBar, State as SearchBarState } from './search-bar/reducer';
import { reducer as survey, State as SurveyState } from './survey/reducer';

export interface RootState {
  bootstrap: BootstrapState;
  cohorts: CohortsState;
  form: FormState;
  matchProfile: MatchProfileState;
  notifications: NotificationsState;
  profile: ProfileState;
  searchBar: SearchBarState;
  userSearch: UserSearchState;
  survey: SurveyState;
};

const reducers = {
  bootstrap,
  cohorts,
  form,
  toast,
  matchProfile,
  notifications,
  profile,
  searchBar,
  userSearch,
  survey,
};

const appReducer: Reducer<RootState> = combineReducers<RootState>(reducers);

export default appReducer;

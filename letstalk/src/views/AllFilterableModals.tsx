import React, { Component } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Alert } from 'react-native';

import { RootState } from '../redux';
import { MultiTrait, MultiTraitTypes } from '../models/multi-trait';
import {
  State as UserSearchState,
  searchByCohort,
  searchByPosition,
  searchBySimpleTrait,
} from '../redux/user-search/reducer';
import {
  ActionTypes as UserSearchActionTypes,
} from '../redux/user-search/actions';
import {
  State as SearchBarState,
  updateFocus,
} from '../redux/search-bar/reducer';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
  SEARCH_LIST_TYPE_CREDENTIALS,
} from '../redux/search-bar/actions';
import { errorToast } from '../redux/toast';
import Colors from '../services/colors';
import { AnalyticsActions, logAnalyticsThenExecuteAsync } from '../services/analytics';
import UserSearchAutocompleteModal from './user-search/UserSearchAutocompleteModal';

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  updateFocus: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  searchByCohort: ActionCreator<
    ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  searchByPosition: ActionCreator<
    ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  searchBySimpleTrait: ActionCreator<
    ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
}

interface Props extends DispatchActions {
  searchBar: SearchBarState;
  onSelectSuccess?(): void;
}

const DEFAULT_SEARCH_SIZE = 10;

class AllFilterableModals extends Component<Props> {

  constructor(props: Props) {
    super(props);

    this.onMultiTraitSelect = this.onMultiTraitSelect.bind(this);
  }

  private async blurSearchBar() {
    await this.props.updateFocus(false);
  }

  private selectSuccess() {
    if (!!this.props.onSelectSuccess) this.props.onSelectSuccess();
  }

  private async searchForTrait(trait: MultiTrait): Promise<void> {
    switch (trait.traitType) {
      case MultiTraitTypes.COHORT:
        this.props.searchByCohort(trait.cohortId, DEFAULT_SEARCH_SIZE);
        break;
      case MultiTraitTypes.POSITION:
        this.props.searchByPosition(trait.roleId, trait.organizationId, DEFAULT_SEARCH_SIZE);
        break;
      case MultiTraitTypes.SIMPLE_TRAIT:
        this.props.searchBySimpleTrait(trait.simpleTraitId, DEFAULT_SEARCH_SIZE);
        break;
      default:
        const _: never = trait;
    }
  }

  private async onMultiTraitSelect(trait: MultiTrait): Promise<void> {
    try {
      logAnalyticsThenExecuteAsync(
        "UserSearch",
        AnalyticsActions.SELECT,
        trait.traitType,
        1,
        this.searchForTrait.bind(this, trait),
      );
      this.selectSuccess();
    } catch (e) {
      await this.props.errorToast(e.errorMsg);
    }
    await this.blurSearchBar();
  }

  render() {
    const {
      value,
      hasFocus,
      suggestions,
      errorMsg,
    } = this.props.searchBar;

    if (!hasFocus) return null;

    return (
      <UserSearchAutocompleteModal
        onSelect={this.onMultiTraitSelect}
        data={suggestions}
        value={value}
        errorMsg={errorMsg}
      />
    );
  }
}


export default connect(({ searchBar }: RootState) => {
  return {
    searchBar,
  };
}, {
  errorToast,
  updateFocus,
  searchByCohort,
  searchByPosition,
  searchBySimpleTrait,
})(AllFilterableModals);

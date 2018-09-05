import React, { Component } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Alert } from 'react-native';

import { RootState } from '../redux';
import { MultiTrait } from '../models/multi-trait';
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
}

interface Props extends DispatchActions {
  searchBar: SearchBarState;
  onSelectSuccess?(): void;
}

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
    // TODO: Actually search here
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
    } = this.props.searchBar;

    if (!hasFocus) return null;

    return (
      <UserSearchAutocompleteModal
        onSelect={this.onMultiTraitSelect}
        data={suggestions}
        value={value}
      />
    );
  }
}


export default connect(({ searchBar }: RootState) => {
  return { searchBar };
}, {
  errorToast,
  updateFocus,
})(AllFilterableModals);

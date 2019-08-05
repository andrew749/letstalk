import * as React from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';

import { RootState } from '../redux';
import { MultiTrait, MultiTraitTypes } from '../models/multi-trait';
import { DEFAULT_SEARCH_SIZE } from '../services/user-search-service';
import {
  State as UserSearchState,
  searchByCohort,
  searchByPosition,
  searchBySimpleTrait,
  searchByGroup,
  setQuery,
} from '../redux/user-search/reducer';
import {
  ActionTypes as UserSearchActionTypes,
  QueryTypes,
} from '../redux/user-search/actions';
import {
  State as SearchBarState,
  updateFocus,
} from '../redux/search-bar/reducer';
import {
  ActionTypes as SearchBarActionTypes,
} from '../redux/search-bar/actions';
import { errorToast } from '../redux/toast';
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
  searchByGroup: ActionCreator<
    ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  setQuery: ActionCreator<
    ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
}

interface StateProps {
  searchBar: SearchBarState;
}

interface ComponentProps {
  onSelectSuccess?(): void;
}

type Props = StateProps & ComponentProps & DispatchActions;

class AllFilterableModals extends React.Component<Props> {

  constructor(props: Props) {
    super(props);

    this.onMultiTraitSelect = this.onMultiTraitSelect.bind(this);
  }

  private async blurSearchBar() {
    await this.props.updateFocus(false);
  }

  private selectSuccess() {
    console.log("Select success");
    if (!!this.props.onSelectSuccess) this.props.onSelectSuccess();
  }

  private async searchForTrait(trait: MultiTrait): Promise<void> {
    switch (trait.traitType) {
      case MultiTraitTypes.COHORT:
        await this.props.setQuery({
          ...trait,
          type: QueryTypes.SEARCH_COHORT,
        });
        this.props.searchByCohort(trait.cohortId, DEFAULT_SEARCH_SIZE);
        break;
      case MultiTraitTypes.POSITION:
        await this.props.setQuery({
          ...trait,
          type: QueryTypes.SEARCH_POSITION,
        });
        this.props.searchByPosition(trait.roleId, trait.organizationId, DEFAULT_SEARCH_SIZE);
        break;
      case MultiTraitTypes.SIMPLE_TRAIT:
        await this.props.setQuery({
          ...trait,
          type: QueryTypes.SEARCH_SIMPLE_TRAIT,
        });
        this.props.searchBySimpleTrait(trait.simpleTraitId, DEFAULT_SEARCH_SIZE);
        break;
      case MultiTraitTypes.GROUP:
        await this.props.setQuery({
          ...trait,
          type: QueryTypes.SEARCH_GROUP,
        });
        this.props.searchByGroup(trait.groupId, DEFAULT_SEARCH_SIZE);
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

interface A{
  errorToast: any
  updateFocus: any
  searchByCohort: any
  searchByPosition: any
  searchBySimpleTrait: any
  searchByGroup: any
  setQuery: any
}


export default connect<StateProps, any, ComponentProps>(({ searchBar }: RootState) => {
  return {
    searchBar,
  };
}, {
  errorToast,
  updateFocus,
  searchByCohort,
  searchByPosition,
  searchBySimpleTrait,
  searchByGroup,
  setQuery,
})(AllFilterableModals as any);
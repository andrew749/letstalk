import Immutable from 'immutable';

import { CohortV2 } from './cohort';
import { GenderId } from './user';

export interface UserSearchResult {
  readonly userId: number;
  readonly firstName: string;
  readonly lastName: string;
  readonly gender: GenderId;
  readonly cohort: CohortV2;
  readonly profilePic: string | null;
  readonly reason: string | null;
}

export interface UserSearchResponse {
  readonly isAnonymous: boolean;
  readonly numResults: number;
  readonly results: Immutable.List<UserSearchResult>;
}

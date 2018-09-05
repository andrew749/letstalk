import Immutable from 'immutable';

import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { UserSearchResponse } from '../models/user-search';
import {
  USER_SEARCH_COHORT_ROUTE,
  USER_SEARCH_POSITION_ROUTE,
  USER_SEARCH_SIMPLE_TRAIT_ROUTE,
} from './constants';

interface BaseUserSearchRequest {
  readonly size: number;
}

interface CohortUserSearchRequest extends BaseUserSearchRequest {
  readonly cohortId: number;
}

interface PositionUserSearchRequest extends BaseUserSearchRequest {
  readonly roleId: number;
  readonly organizationId: number;
}

interface SimpleTraitUserSearchRequest extends BaseUserSearchRequest {
  readonly simpleTraitId: number;
}

type UserSearchRequest =
  | CohortUserSearchRequest
  | PositionUserSearchRequest
  | SimpleTraitUserSearchRequest

class UserSearchService {
  private requestor: Requestor
  private auth: Auth

  constructor(requestor: Requestor, auth: Auth) {
    this.requestor = requestor;
    this.auth = auth;
  }

  private async doUserSearch(url: string, req: UserSearchRequest): Promise<UserSearchResponse> {
    const sessionToken = await this.auth.getSessionToken();
    const res: UserSearchResponse = await this.requestor.post(url, req, sessionToken);
    return {
      ...res,
      results: Immutable.List(res.results),
    };
  }

  async searchByCohort(req: CohortUserSearchRequest): Promise<UserSearchResponse> {
    return this.doUserSearch(USER_SEARCH_COHORT_ROUTE, req);
  }

  async searchByPosition(req: PositionUserSearchRequest): Promise<UserSearchResponse> {
    return this.doUserSearch(USER_SEARCH_POSITION_ROUTE, req);
  }

  async searchBySimpleTrait(req: SimpleTraitUserSearchRequest): Promise<UserSearchResponse> {
    return this.doUserSearch(USER_SEARCH_SIMPLE_TRAIT_ROUTE, req);
  }
}

const userSearchService = new UserSearchService(requestor, auth);
export default userSearchService;

import {apiModule, pluggableApiModule} from './api_module';

// This needs to be unique from other apis
export const API_NAME = "fetchMatchingRoundsApi";
// this needs to match the name of the function in HiveApiService
export const API_FUNC = "getMatchingRounds";

const initialDataState = {
    groupId: undefined,
}

export const fetchMatchingRoundsApiModule = apiModule(API_NAME);
export const fetchMatchingRoundsApi = pluggableApiModule(
    fetchMatchingRoundsApiModule,
    API_FUNC, 
    initialDataState,
);
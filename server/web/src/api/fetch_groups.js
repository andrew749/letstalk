import {apiModule, pluggableApiModule} from './api_module';

// This needs to be unique from other apis
export const API_NAME = "fetchGroupsApi";
// this needs to match the name of the function in HiveApiService
export const API_FUNC = "fetchGroups";

const initialDataState = {
}

export const fetchGroupsApiModule = apiModule(API_NAME);
export const fetchGroupsApi = pluggableApiModule(
    fetchGroupsApiModule,
    API_FUNC, 
    initialDataState,
);
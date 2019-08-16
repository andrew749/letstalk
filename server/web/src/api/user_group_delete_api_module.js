import {apiModule, pluggableApiModule} from './api_module';

export const API_NAME = "userGroupDeleteApi";
export const API_FUNC = "deleteMemberFromGroup";

const initialDataState = {
    userId: undefined,
    groupId: undefined,
}

export const userGroupDeleteApiModule = apiModule(API_NAME);
export const userGroupDeleteApi = pluggableApiModule(
    userGroupDeleteApiModule,
    API_FUNC, 
    initialDataState,
);
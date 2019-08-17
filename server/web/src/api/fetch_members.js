import {apiModule, pluggableApiModule} from './api_module';
import {fetchGroupsApiModule} from './fetch_groups';

// This needs to be unique from other apis
export const API_NAME = "fetchMembersApi";
// this needs to match the name of the function in HiveApiService
export const API_FUNC = "fetchMembers";

const initialDataState = {
    groupId: undefined,
}

const rawApiModule = apiModule(API_NAME);
export const fetchMembersApiModule = { 
    getCurrentGroup: (state) => {
        if (rawApiModule.getParams(state)) {
            if (fetchGroupsApiModule.getData(state)) {
                const matchingGroups = fetchGroupsApiModule.getData(state).managedGroups;
                if (matchingGroups) {
                    const filteredGroups = matchingGroups.filter(group => group.groupId == rawApiModule.getParams(state).groupId);
                    if (filteredGroups.length > 0) {
                        return filteredGroups[0];
                    }
                }
            }
        }
        return undefined;
    },
    ...rawApiModule,
};
export const fetchMembersApi = pluggableApiModule(
    fetchMembersApiModule,
    API_FUNC, 
    initialDataState,
);
import {apiModule, pluggableApiModule} from './api_module';
import {HiveApiService} from './api_controller';

export const API_NAME = "matchRoundApi";
export const API_FUNC = "createNewMatchingRound";

const initialDataState = {
    groupId: undefined,
    userIds: undefined,
    maxLowerYearsPerUpperYear: undefined, 
    maxUpperYearsPerLowerYear: undefined, 
    youngestUpperGradYear: undefined,
}

export const matchRoundApiModule = apiModule(API_NAME);
export const matchRoundApi = pluggableApiModule(
    matchRoundApiModule,
    API_FUNC, 
    initialDataState,
);

const deleteMatchRoundInitialDataState = {
    matchRoundId: undefined,
}

export const DELETE_API_NAME = "deleteMatchRoundApi";
export const DELETE_API_FUNC = "deleteMatchingRound";

export const deleteMatchRoundApiModule = apiModule(DELETE_API_NAME);
export const deleteMatchRoundApi = pluggableApiModule(
    deleteMatchRoundApiModule,
    DELETE_API_FUNC,
    deleteMatchRoundInitialDataState,
);

const commitMatchRoundInitialDataState = {
    matchRoundId: undefined
};

export const COMMIT_MATCH_ROUND_API_NAME = "commitMatchRoundApi";
export const COMMIT_MATCH_ROUND_API_FUNC = "commitMatchingRound";

export const commitMatchRoundApiModule = apiModule(COMMIT_MATCH_ROUND_API_NAME);
export const commitMatchRoundApi = pluggableApiModule(
    commitMatchRoundApiModule,
    COMMIT_MATCH_ROUND_API_FUNC,
    commitMatchRoundInitialDataState,
);
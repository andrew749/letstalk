import {apiModule, pluggableApiModule} from './api_module';
import {HiveApiService} from './api_controller';

export const API_NAME = "matchRoundApi";
export const API_FUNC = "createNewMatchingRound";

const initialDataState = {
    groupId: undefined,
    userIds: undefined,
    maxLowerYearsPerUpperYear: undefined, 
    maxUpperYearsPerLowerYear: undefined, 
    youngestUpperYearGrad: undefined,
}

export const matchRoundApiModule = apiModule(API_NAME);
export const matchRoundApi = pluggableApiModule(
    matchRoundApiModule,
    API_FUNC, 
    initialDataState,
);
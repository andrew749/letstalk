import {apiModule, pluggableApiModule} from './api_module';
import {HiveApiService} from './api_controller';

export const API_NAME = "getCohortsApi";
export const API_FUNC = "getCohorts";

const initialDataState = {
    cohorts: undefined,
}

export const getCohortsApiModule = apiModule(API_NAME);
export const getCohortsApi = pluggableApiModule(
    getCohortsApiModule,
    API_FUNC,
    initialDataState,
);

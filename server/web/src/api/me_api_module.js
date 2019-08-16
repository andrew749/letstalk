import {apiModule, pluggableApiModule} from './api_module';

// This needs to be unique from other apis
export const API_NAME = "meApi";
// this needs to match the name of the function in HiveApiService
export const API_FUNC = "me";

const initialDataState = {
    /* This api has no data */
}

export const meApiModule = apiModule(API_NAME);
export const meApi = pluggableApiModule(
    meApiModule,
    API_FUNC, 
    initialDataState,
);
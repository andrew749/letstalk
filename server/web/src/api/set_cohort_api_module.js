import {apiModule, pluggableApiModule} from './api_module';

export const API_NAME = "setCohortApi";
export const API_FUNC = "setCohort";

const initialDataState = {};

export const setCohortApiModule = apiModule(API_NAME);
export const setCohortApi = pluggableApiModule(
    setCohortApiModule,
    API_FUNC,
    initialDataState,
);

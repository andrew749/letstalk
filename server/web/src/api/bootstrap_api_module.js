import {apiModule, pluggableApiModule} from './api_module';
import {HiveApiService} from './api_controller';

export const API_NAME = "bootstrapApi";
export const API_FUNC = "getBootstrap";

const initialDataState = {
    state: undefined,
    connections: undefined,
}

export const bootstrapApiModule = apiModule(API_NAME);
export const bootstrapApi = pluggableApiModule(
    bootstrapApiModule,
    API_FUNC,
    initialDataState,
);

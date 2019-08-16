import {apiModule, pluggableApiModule} from './api_module';

export const API_NAME = "sendVerificationEmailApi";
export const API_FUNC = "sendVerificationEmail";

const initialDataState = {}

export const sendVerificationEmailApiModule = apiModule(API_NAME);
export const sendVerificationEmailApi = pluggableApiModule(
    sendVerificationEmailApiModule,
    API_FUNC,
    initialDataState,
);

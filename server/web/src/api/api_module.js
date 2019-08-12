
const EXECUTE_API_CALL = (apiName) => `${apiName}_EXECUTE_API_CALL`;
const EXECUTING_API_CALL = (apiName) => `${apiName}_EXECUTING_API_CALL`;
const EXECUTED_API_CALL_SUCCESSFULLY = (apiName) => `${apiName}_EXECUTED_API_CALL_SUCCESSFULLY`;
const ERROR_EXECUTING_API_CALL = (apiName) => `${apiName}_ERROR_EXECUTING_API_CALL`;

const initialApiState = (initialDataState) => ({
    shouldFetch: false,
    doneFetching: false,
    hasError: false,
    data: initialDataState,
    apiError: undefined,
});

export function executeApiCallAction(apiName, params) {
    return {type: EXECUTE_API_CALL(apiName), params: params};
}

export function executingApiCallAction(apiName) {
    return {type: EXECUTING_API_CALL(apiName)};
}

export function executeApiCallSuccessfullyAction(apiName, result) {
    return {type: EXECUTED_API_CALL_SUCCESSFULLY(apiName), data: result};
}

export function errorExecutingApiCallAction(apiName, err) {
    return {type: ERROR_EXECUTING_API_CALL(apiName), apiError: err};
}

export const apiModule = (apiName) => ({
    bindApi: (apiFunction) => (state, dispatch) => (params) => 
        HiveApiService(state, dispatch)[apiFunction]({
            ...params, 
            started: () => dispatch(executingApiCallAction(apiName)),
            done: (data) => dispatch(executeApiCallSuccessfullyAction(apiName, data.Result)),
            error: (err) => dispatch(errorExecutingApiCallAction(apiName, err))
        }),
    shouldExecuteApiCall: (state) => {
        console.log(state);
        return !!state[apiName].shouldFetch;
    },
    bindReducer: (initialDataState) => (state = initialApiState(initialDataState), action) => {
        switch (action.type) {
            case EXECUTE_API_CALL:
                return Object.assign({}, state, { shouldFetch: true, doneFetching: false, params: action.params });
            case EXECUTING_API_CALL:
                return Object.assign({}, state, { shouldFetch: false });
            case EXECUTED_API_CALL_SUCCESSFULLY:
                return Object.assign({}, state, { shouldFetch: false, doneFetching: true, hasError: false, data: action.data });
            case ERROR_EXECUTING_API_CALL:
                return Object.assign({}, state, { shouldFetch: false, doneFetching: true, hasError: true, apiError: action.apiError });
            default:
                return state;
        }
    },
    getParams: (state) => {
        return state[apiName].params;
    },
    isFinished: (state) => {
        return state[apiName].doneFetching;
    },
    getData: (state) => {
        return state[apiName].data;
    },
    hasError: (state) =>  {
        return state[apiName].hasError;
    },
    getErrorMessage: (state) => {
        return state[apiName].apiError;
    },
    getApiExecuteAction: () => {
        return executeApiCallAction(apiName);
    },
});

export const pluggableApiModule = (module, apiFunc, initialDataState) => ({
    call: (params, state, dispatch) => module.bindApi(apiFunc)(state, dispatch)(params),
    reducer: module.bindReducer(initialDataState),
    module: {...module},
});
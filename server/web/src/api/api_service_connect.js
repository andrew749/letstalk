import { connect } from 'react-redux';
import {HiveApiService} from './api_controller';

/**
 * Specialized connect that injects an api service into the component 
 * @param {*} mapStateToProps 
 * @param {*} mapDispatchToProps 
 * @param {*} mergeProps 
 */
export default function(mapStateToProps, mapDispatchToProps, mergeProps) {
    return connect(
        (state, ownProps) => ({
            ...functionOrEmptyObject(mapStateToProps, state, ownProps),
            ...ownProps,
            apiServiceReducer: state.apiServiceReducer,
        }),
        (dispatch, ownProps) => ({
            ...functionOrEmptyObject(mapDispatchToProps, dispatch, ownProps),
            ...ownProps,
            dispatch: dispatch,
        }),
        (stateProps, dispatchProps, ownProps) => ({
             ...functionOrEmptyObject(mergeProps, stateProps, dispatchProps, ownProps),
             ...stateProps,
             ...dispatchProps,
             ...ownProps,
            apiService: HiveApiService({...stateProps, ...ownProps}, dispatchProps.dispatch),
        }),
    );
}

function functionOrEmptyObject(f, ...args) {
    let res = {};
    if (!!f) {
        res = f(...args);
    }
    return res;
}
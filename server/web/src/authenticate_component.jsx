import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import apiServiceConnect from './api/api_service_connect';

const SET_REDIRECT_ROUTE_ACTION = 'POST_AUTH/SET_REDIRECT_ROUTE_ACTION';

const initialState = {
  redirectPath: null,
}

function getSetRedirectRoute(redirectPath) {
    return {type: SET_REDIRECT_ROUTE_ACTION, redirectPath: redirectPath};
}

export function postAuthReducer(state = initialState, action) {
    switch(action.type) {
        case SET_REDIRECT_ROUTE_ACTION:
            return Object.assign({}, state, {redirectPath: action.redirectPath});
        default:
            return state;
    }
}

/**
 * A specialized react-router Route that can redirect the user to a login page
 * if they are not authenticated.
 *
 *  loginPath: the path to redirect a user to if they are unauthenticated
 */
const AuthenticatedRoute = ({ component: Component, ...rest }) => {
      if (!rest.isAuthenticated) {
          rest.setRedirectRoute(rest.location.pathname);
      }
      return (
          <Route
              render={(props) => {
                  if (!!rest.isAuthenticated) {
                      return <Component {...props} />;
                  } else {
                      return <Redirect to={{
                          pathname: rest.loginPath,
                          state: { from: props.location }
                      }}/>;
                  }
              }}
              {...rest}
          />
      );
}

function isAuthenticated(state) {
    return state.apiServiceReducer.sessionId && state.apiServiceReducer.isValid;
}

const AuthenticatedComponent = apiServiceConnect(
    (state) => {
        return {
            isAuthenticated: isAuthenticated(state)
        };
    },
    (dispatch) => {
        return {
            setRedirectRoute: (route) => dispatch(getSetRedirectRoute(route))
        };
    }
)(AuthenticatedRoute);

export default AuthenticatedComponent;


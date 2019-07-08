import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';

/**
 * A specialized react-router Route that can redirect the user to a login page
 * if they are not authenticated.
 * 
 *  loginPath: the path to redirect a user to if they are unauthenticated
 */
const AuthenticatedRoute = ({ component: Component, ...rest }) =>
<Route
  {...rest}
  render={(props) => {
      return !!rest.isAuthenticated
      ? <Component {...props} />
      : <Redirect to={{
          pathname: rest.loginPath,
          state: { from: props.location }
      }}/>
  }
  }
/>;

function isAuthenticated(state) {
    return state.loginReducer.isAuthenticated;
}

const AuthenticatedComponent = connect(
    (state) => {
        return {
            isAuthenticated: isAuthenticated(state)
        };
    },
    null
)(AuthenticatedRoute);

export default AuthenticatedComponent;


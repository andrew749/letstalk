import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { loginPath } from './routes.js';
import { connect } from 'react-redux';
import {withCookies} from 'react-cookie';

import CookieAwareComponent from './cookie_aware_component.jsx';

const AuthenticatedRoute = ({ component: Component, ...rest }) =>
<Route
  {...rest}
  render={(props) => {
      return !!rest.isAuthenticated
      ? <Component {...props} />
      : <Redirect to={{
          pathname: loginPath,
          state: { from: props.location }
      }}/>

  }
  }
/>;

const AuthenticatedComponent = connect(
    (state) => {
        return {
            isAuthenticated: state.loginReducer.isAuthenticated
        };
    },
    null
)(AuthenticatedRoute);

export default AuthenticatedComponent;


import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, Redirect, BrowserRouter} from 'react-router-dom';
import {withCookies} from 'react-cookie';

import CookieAwareComponent from '../cookie_aware_component.jsx';
import SignupPage, {signupReducer} from './signup.jsx';
import LoginPage, {loginReducer} from './login.jsx';
import LandingPage from '../landing.jsx';

import AuthenticatedRoute from '../authenticate_component.jsx';
import { loginPathWeb, signupPathWeb, landingPathWeb } from '../routes.js';
import HiveHeader from './components/hive_header.jsx';

const reducers = combineReducers({
    signupReducer,
    loginReducer,
});

const store = createStore(reducers);

// Specialized AuthenticatedRoute component for the normal login page.
const AuthenticatedRouteWebapp = (...props) => 
    <AuthenticatedRoute 
        {...props} 
        loginPath={loginPathWeb}
        />;

class App extends React.Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <HiveHeader />
                    <Switch>
                        <Route path={loginPathWeb} component={LoginPage} />
                        <Route path={signupPathWeb} component={SignupPage} />
                        <AuthenticatedRouteWebapp exact path={landingPathWeb} component={LandingPage} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

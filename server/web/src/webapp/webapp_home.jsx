import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, BrowserRouter} from 'react-router-dom';
import {withCookies} from 'react-cookie';

import CookieAwareComponent from '../cookie_aware_component.jsx';
import SignupPage, {signupReducer} from '../signup.jsx';
import LoginPage, {loginReducer} from '../login.jsx';
import LandingPage from './landing.jsx';
import GroupRegisterPage, {groupRegisterReducer} from './group_register_page.jsx';
import Cookies from 'universal-cookie';

import AuthenticatedRoute from '../authenticate_component.jsx';
import { loginPathWeb, signupPathWeb, landingPathWeb, registerWithGroupPathWeb } from '../routes.js';
import HiveHeader from './components/hive_header.jsx';
import {apiServiceReducer, HiveApiService} from '../api/api_controller';

const reducers = combineReducers({
    apiServiceReducer,
    signupReducer,
    loginReducer,
    groupRegisterReducer,
});

const store = createStore(reducers);

function onLoad() {
    let sessionId = (new Cookies()).get("sessionId");

    if (!!sessionId) {
        console.log("Loaded cookie: " + sessionId);
        HiveApiService(store.getState(), store.dispatch).setSessionId(sessionId);
    }
}

onLoad();

// Specialized AuthenticatedRoute component for the normal login page.
const AuthenticatedRouteWebapp = (props) => {
    return <AuthenticatedRoute
        loginPath={loginPathWeb}
        {...props}
        />;
}

class App extends React.Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <HiveHeader />
                    <Switch>
                        <Route path={loginPathWeb} component={LoginPage} />
                        <Route path={signupPathWeb} component={SignupPage} isAdminPage={false} />
                        <AuthenticatedRouteWebapp path={registerWithGroupPathWeb} component={GroupRegisterPage} />
                        <AuthenticatedRouteWebapp exact path={landingPathWeb} component={LandingPage} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

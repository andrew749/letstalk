import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, BrowserRouter} from 'react-router-dom';
import {withCookies} from 'react-cookie';

import CookieAwareComponent from '../cookie_aware_component.jsx';
import SignupPage from '../signup.jsx';
import LoginPage, {loginReducer} from '../login.jsx';
import LandingPage from './landing.jsx';
import GroupRegisterPage, {groupRegisterReducer} from './group_register_page.jsx';
import VerifyEmailPage from './verify_email';
import Cookies from 'universal-cookie';

import AuthenticatedRoute from '../authenticate_component.jsx';
import {
    loginPathWeb,
    signupPathWeb,
    landingPathWeb,
    registerWithGroupPathWeb,
    verifyEmailPathWeb,
    setCohortPathWeb,
    surveyPathWeb,
} from '../routes.js';
import HiveHeader from './components/hive_header.jsx';
import {apiServiceReducer, HiveApiService} from '../api/api_controller';
import {API_NAME as ME_API, meApi} from '../api/me_api_module';
import {API_NAME as BOOTSTRAP_API, bootstrapApi} from '../api/bootstrap_api_module';
import {
    API_NAME as SEND_VERIFICATION_EMAIL,
    sendVerificationEmailApi,
} from '../api/send_verification_email_api_module';
import {postAuthReducer} from '../authenticate_component';

const apiModules = {
    [ME_API]: meApi,
    [BOOTSTRAP_API]: bootstrapApi,
    [SEND_VERIFICATION_EMAIL]: sendVerificationEmailApi,
}

// build reducer dict
let apiModuleReducers = {};
Object.keys(apiModules).forEach((key) => apiModuleReducers[key] = apiModules[key].reducer);

const reducers = combineReducers({
    apiServiceReducer,
    loginReducer,
    groupRegisterReducer,
    postAuthReducer,
    ...apiModuleReducers,
});

const DEBUG = false;
const store = createStore(reducers);

function consoleDebug(s) {
  if (DEBUG) console.log(s);
}

function onLoad() {
    let sessionId = (new Cookies()).get("sessionId");

    store.subscribe(() => {
        let state = store.getState();
        Object.keys(apiModules).forEach( (key) => {
            consoleDebug(`Evaluting api module ${key}`);
            let mod = apiModules[key];
            if (!!mod.module.shouldExecuteApiCall(state)) {
                consoleDebug(`Executing api call for  api module ${key}`);
                let params = mod.module.getParams(state);
                mod.call(params, state, store.dispatch);
            }
            // pull state again since `store.subscribe` works recursively (if some state change
            // happens in this closure, then it will call subscribe again and finish this up once
            // that is done. Since the state can change in the middle, it is possible that on the
            // next iteration, the state is outdated, causing duplicate calls.
            state = store.getState();
        });
    });

    if (!!sessionId) {
        consoleDebug("Loaded cookie: " + sessionId);
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
                        <Route path={loginPathWeb} component={LoginPage} isAdminApp={false} />
                        <Route path={signupPathWeb} component={SignupPage} isAdminApp={false} />
                        <AuthenticatedRouteWebapp exact path={landingPathWeb} component={LandingPage} />
                        <AuthenticatedRouteWebapp path={registerWithGroupPathWeb} component={GroupRegisterPage} />
                        <AuthenticatedRouteWebapp path={verifyEmailPathWeb} component={VerifyEmailPage} />
                        <AuthenticatedRouteWebapp path={setCohortPathWeb} component={GroupRegisterPage} />
                        <AuthenticatedRouteWebapp path={surveyPathWeb} component={GroupRegisterPage} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

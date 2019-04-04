import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, Redirect, BrowserRouter} from 'react-router-dom';
import {withCookies} from 'react-cookie';

import CookieAwareComponent from './cookie_aware_component.jsx';
import LoginPage, {loginReducer} from './login.jsx';
import AdhocAddPage from './adhoc_add.jsx';
import LandingPage from './landing.jsx';

import AuthenticatedRoute from './authenticate_component.jsx';
import { URL_PREFIX, loginPath, adhocAddToolPath, landingPath } from './routes.js';
import HiveToolTemplate from './hive_tool_template.jsx';

const reducers = combineReducers({
    loginReducer
});

const store = createStore(reducers);

class App extends React.Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <HiveToolTemplate />
                    <Switch>
                        <Route path={loginPath} component={LoginPage} />
                        <AuthenticatedRoute exact path={landingPath} component={LandingPage} />
                        <AuthenticatedRoute path={adhocAddToolPath} component={AdhocAddPage} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

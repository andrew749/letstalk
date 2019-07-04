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
import DeleteUserToolPage from './user_delete_tool.jsx';

import AuthenticatedRoute from './authenticate_component.jsx';
import { URL_PREFIX, loginPath, adhocAddToolPath, landingPath, deleteUserToolPath } from './routes.js';
import HiveToolTemplate from './hive_tool_template.jsx';

const reducers = combineReducers({
    loginReducer
});

const store = createStore(reducers);

// Specialize the general AuthenticatedRoute component to work with admin login page.
const AuthenticatedRouteAdmin = (...props) => 
    <AuthenticatedRoute 
        loginPath={loginPath} 
        {...props}
        />;

class App extends React.Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <HiveToolTemplate />
                    <Switch>
                        <Route path={loginPath} component={LoginPage} />
                        <AuthenticatedRouteAdmin exact path={landingPath} component={LandingPage} />
                        <AuthenticatedRouteAdmin path={adhocAddToolPath} component={AdhocAddPage} />
                        <AuthenticatedRouteAdmin path={deleteUserToolPath} component={DeleteUserToolPage} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

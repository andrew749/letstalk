import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, BrowserRouter} from 'react-router-dom';
import {withCookies} from 'react-cookie';

import CookieAwareComponent from './cookie_aware_component.jsx';
import LoginPage, {loginReducer} from './login.jsx';
import SignupPage, {signupReducer} from './signup.jsx';
import AdhocAddPage from './adhoc_add.jsx';
import LandingPage from './landing.jsx';
import MatchingPage from './matching';
import MembersPage from './members';
import DeleteUserToolPage from './user_delete_tool.jsx';
import ManagedGroupPage from './managed_group.jsx';
import {getManagedGroupsReducer, getShouldFetchGroups, fetchingGroupsAction, gotGroupsAction, errorFetchingGroupsAction} from './get_managed_groups_view'
import {HiveApiService} from './api_controller.js';

import AuthenticatedRoute from './authenticate_component.jsx';
import { loginPath, signupPath, adhocAddToolPath, landingPath, deleteUserToolPath, groupManagementToolPath, matchingPath, membersPath } from './routes.js';
import HiveToolTemplate from './hive_tool_template.jsx';

const reducers = combineReducers({
    loginReducer,
    signupReducer,
    getManagedGroupsReducer
});

const store = createStore(reducers);

function getSessionIdFromState(state) {
    return state.loginReducer.sessionId;
}

store.subscribe(() => {
    let sessionId = getSessionIdFromState(store.getState());
    // set the session id after logging
    if (!!sessionId) {
        HiveApiService.setSessionId(sessionId);
    }

    // if somebody posted a fetch event, then get the api
    let shouldFetchGroups = getShouldFetchGroups(store.getState());
    if (!!shouldFetchGroups) {
        HiveApiService.fetchGroups(
            () => {store.dispatch(fetchingGroupsAction())},
            (data) => {store.dispatch(gotGroupsAction(data.Result.managedGroups))},
            (err) => {store.dispatch(errorFetchingGroupsAction(err))}
        );
    }
})


// Specialize the general AuthenticatedRoute component to work with admin login page.
const AuthenticatedRouteAdmin = (props) =>
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
                        <Route path={loginPath} render={(props) => <LoginPage {...props} isAdminPage={true} />} />
                        <Route path={signupPath} render={(props) => <SignupPage {...props} isAdminPage={true} />}  />
                        <AuthenticatedRouteAdmin exact path={landingPath} component={LandingPage} />
                        <AuthenticatedRouteAdmin path={membersPath} component={MembersPage} />
                        <AuthenticatedRouteAdmin path={matchingPath} component={MatchingPage} />
                        <AuthenticatedRouteAdmin path={adhocAddToolPath} component={AdhocAddPage} />
                        <AuthenticatedRouteAdmin path={deleteUserToolPath} component={DeleteUserToolPage} />
                        <AuthenticatedRouteAdmin path={groupManagementToolPath} component={ManagedGroupPage} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

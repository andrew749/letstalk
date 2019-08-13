import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, BrowserRouter} from 'react-router-dom';
import {withCookies, CookiesProvider} from 'react-cookie';
import Cookies from 'universal-cookie';

import CookieAwareComponent from './cookie_aware_component.jsx';
import LoginPage, {loginReducer} from './login.jsx';
import SignupPage, {signupReducer} from './signup.jsx';
import ModalContainer, {modalReducer} from './modal_container.jsx';
import AdhocAddPage from './adhoc_add.jsx';
import LandingPage from './landing.jsx';
import MatchingPage, {matchingReducer, getShouldFetchMatchingRoundsForGroup, getMatchingRoundsGroupToFetch, fetchingMatchingRoundsForGroupAction, fetchedMatchingRoundsForGroupAction, errorFetchingMatchingRoundsForGroupAction} from './matching';
import MembersPage from './members';
import DeleteUserToolPage from './user_delete_tool.jsx';
import ManagedGroupPage from './managed_group.jsx';
import {API_NAME as MATCH_ROUND_API, matchRoundApi, DELETE_API_NAME as DELETE_MATCH_ROUND_API_NAME, deleteMatchRoundApi} from './api/match_round_api_module';
import {API_NAME as DELETE_USER_GROUP_API, userGroupDeleteApi} from './api/user_group_delete_api_module';
import {getManagedGroupsReducer, getShouldFetchGroups, fetchingGroupsAction, gotGroupsAction, errorFetchingGroupsAction} from './get_managed_groups_view'
import {membersReducer, getShouldFetchMembers, fetchingMembersAction, gotMembersAction, errorFetchingMembersAction, getGroupToFetch} from './members';
import {apiServiceReducer, HiveApiService, getShouldFetchProfile, didFetchProfileAction, fetchingProfileAction, fetchProfileErrorAction} from './api/api_controller';

import AuthenticatedRoute from './authenticate_component.jsx';
import { loginPath, signupPath, adhocAddToolPath, landingPath, deleteUserToolPath, groupManagementToolPath, matchingPath, membersPath } from './routes.js';
import NavbarContainer from './navbar_container.jsx';

const apiModules = {
    [MATCH_ROUND_API]: matchRoundApi,
    [DELETE_USER_GROUP_API]: userGroupDeleteApi,
    [DELETE_MATCH_ROUND_API_NAME]: deleteMatchRoundApi,
}

// build reducer dict
let apiModuleReducers = {};
console.log(apiModules)
Object.keys(apiModules).forEach((key) => apiModuleReducers[key] = apiModules[key].reducer);
console.log(apiModuleReducers);

const reducers = combineReducers({
    apiServiceReducer,
    loginReducer,
    signupReducer,
    getManagedGroupsReducer,
    membersReducer,
    modalReducer,
    matchingReducer,
    ...apiModuleReducers,
});

const store = createStore(reducers);

function onLoad() {
    let sessionId = (new Cookies()).get("sessionId");

    if (!!sessionId) {
        console.log("Loaded cookie: " + sessionId);
        HiveApiService(store.getState(), store.dispatch).setSessionId(sessionId);
    }

    store.subscribe(() => {
        // if somebody posted a fetch event, then get the api
        let shouldFetchGroups = getShouldFetchGroups(store.getState());
        if (!!shouldFetchGroups) {
            HiveApiService(store.getState(), store.dispatch).fetchGroups(
                () => {store.dispatch(fetchingGroupsAction())},
                (data) => {store.dispatch(gotGroupsAction(data.Result.managedGroups))},
                (err) => {store.dispatch(errorFetchingGroupsAction(err))}
            );
        }

        let shouldFetchProfile = getShouldFetchProfile(store.getState());
        if (!!shouldFetchProfile) {
            console.log("Fetching profile");
            HiveApiService(store.getState(), store.dispatch).me(
                () => {store.dispatch(fetchingProfileAction())},
                (data) => {store.dispatch(didFetchProfileAction(data.Result))},
                (err) => {store.dispatch(fetchProfileErrorAction(err))}
            );
        }

        let shouldFetchMatchingRounds = getShouldFetchMatchingRoundsForGroup(store.getState());
        if (!!shouldFetchMatchingRounds) {
            let matchingRoundsGroupToFetch = getMatchingRoundsGroupToFetch(store.getState());
            console.log("Fetching matching rounds for " + matchingRoundsGroupToFetch);
            HiveApiService(store.getState(), store.dispatch).getMatchingRounds(
                matchingRoundsGroupToFetch.groupId,
                () => {store.dispatch(fetchingMatchingRoundsForGroupAction())},
                (data) => {store.dispatch(fetchedMatchingRoundsForGroupAction(data.Result))},
                (err) => {store.dispatch(errorFetchingMatchingRoundsForGroupAction(err))}
            );
        }
        let state = store.getState();

        Object.keys(apiModules).forEach( (key) => {
            console.log(`Evaluting api module ${key}`);
            let mod = apiModules[key];
            if (!!mod.module.shouldExecuteApiCall(state)) {
                console.log(`Executing api call for  api module ${key}`);
                let params = mod.module.getParams(state);
                mod.call(params, state, store.dispatch);
            }
        });

        // TODO: Finish this part, write the routes, view results???
        let shouldFetchMembers = getShouldFetchMembers(store.getState());
        let groupToFetch = getGroupToFetch(store.getState());
        if (!!shouldFetchMembers && !!groupToFetch) {
            console.log("Fetching group " + groupToFetch.groupName);
            HiveApiService(store.getState(), store.dispatch).fetchMembers(
                groupToFetch.groupId,
                () => {store.dispatch(fetchingMembersAction())},
                (data) => {store.dispatch(gotMembersAction(data.Result))},
                (err) => {store.dispatch(errorFetchingMembersAction(err))}
            );
        }
    });
}

onLoad();

// Specialize the general AuthenticatedRoute component to work with admin login page.
const AuthenticatedRouteAdmin = (props) =>
    <AuthenticatedRoute
        loginPath={loginPath}
        {...props}
        />;

class App extends React.Component {
    
    render() {
        return (
            <CookiesProvider>
                <Provider store={store}>
                    <ModalContainer/>
                    <BrowserRouter>
                        <NavbarContainer />
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
            </CookiesProvider>
        );
    }
}

const CookieAwarePage = CookieAwareComponent(withCookies(App));
ReactDOM.render(<CookieAwarePage/>, document.getElementById('content'));

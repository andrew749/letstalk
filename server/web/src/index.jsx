import ReactDOM from 'react-dom';
import React from 'react';
import { createStore, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import { Switch, Route, BrowserRouter} from 'react-router-dom';
import {withCookies, CookiesProvider} from 'react-cookie';
import Cookies from 'universal-cookie';

import CookieAwareComponent from './cookie_aware_component.jsx';
import LoginPage, {loginReducer} from './login.jsx';
import SignupPage from './signup.jsx';
import ModalContainer, {modalReducer} from './modal_container.jsx';
import AdhocAddPage from './adhoc_add.jsx';
import LandingPage from './landing.jsx';
import MatchingPage, {matchingReducer, shouldFetchMatchingRoundsForGroupAction, getShouldFetchMatchingRoundsForGroup, fetchingMatchingRoundsForGroupAction, fetchedMatchingRoundsForGroupAction, errorFetchingMatchingRoundsForGroupAction} from './matching';
import MembersPage from './members';
import DeleteUserToolPage from './user_delete_tool.jsx';
import ManagedGroupPage from './managed_group.jsx';
import {API_NAME as MATCH_ROUND_API, matchRoundApi, DELETE_API_NAME as DELETE_MATCH_ROUND_API_NAME, deleteMatchRoundApi, COMMIT_MATCH_ROUND_API_NAME, commitMatchRoundApi} from './api/match_round_api_module';
import {API_NAME as DELETE_USER_GROUP_API, userGroupDeleteApi} from './api/user_group_delete_api_module';
import {API_NAME as ME_API, meApi} from './api/me_api_module';
import {API_NAME as FETCH_GROUPS_API, fetchGroupsApi} from './api/fetch_groups';
import {API_NAME as FETCH_MEMBERS_API, fetchMembersApi, fetchMembersApiModule} from './api/fetch_members';
import {apiServiceReducer, HiveApiService} from './api/api_controller';
import {GroupContextManager, groupContextReducer, getCurrentGroup} from './group_context_reducer';

import AuthenticatedRoute, {postAuthReducer} from './authenticate_component.jsx';
import { loginPath, signupPath, adhocAddToolPath, landingPath, deleteUserToolPath, groupManagementToolPath, matchingPath, membersPath } from './routes.js';
import NavbarContainer from './navbar_container.jsx';

import 'react-bootstrap-table-next/dist/react-bootstrap-table2.min.css';

const apiModules = {
    [MATCH_ROUND_API]: matchRoundApi,
    [DELETE_USER_GROUP_API]: userGroupDeleteApi,
    [DELETE_MATCH_ROUND_API_NAME]: deleteMatchRoundApi,
    [ME_API]: meApi,
    [COMMIT_MATCH_ROUND_API_NAME]: commitMatchRoundApi,
    [FETCH_GROUPS_API]: fetchGroupsApi,
    [FETCH_MEMBERS_API]: fetchMembersApi,
}

// build reducer dict
let apiModuleReducers = {};
console.log(apiModules)
Object.keys(apiModules).forEach((key) => apiModuleReducers[key] = apiModules[key].reducer);
console.log(apiModuleReducers);

const reducers = combineReducers({
    apiServiceReducer,
    loginReducer,
    modalReducer,
    matchingReducer,
    postAuthReducer,
    groupContextReducer,
    ...apiModuleReducers,
});

const store = createStore(reducers);
let groupContextManager = new GroupContextManager();
groupContextManager.subscribeListenerToNewGroupAction((state, group) => group && store.dispatch(fetchMembersApiModule.getApiExecuteAction({groupId: group.groupId})));
groupContextManager.subscribeListenerToNewGroupAction((state, group) => group && store.dispatch(shouldFetchMatchingRoundsForGroupAction(group)));

function onLoad() {
    let sessionId = (new Cookies()).get("sessionId");

    if (!!sessionId) {
        console.log("Loaded cookie: " + sessionId);
        HiveApiService(store.getState(), store.dispatch).setSessionId(sessionId);
    }

    store.subscribe(() => {
        let shouldFetchMatchingRounds = getShouldFetchMatchingRoundsForGroup(store.getState());
        if (!!shouldFetchMatchingRounds) {
            let matchingRoundsGroupToFetch = getCurrentGroup(store.getState());
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
        // check if we need to update the global state based on a group change action
        groupContextManager.checkGroupDidChange(state);
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
                            <Route path={loginPath} render={(props) => <LoginPage {...props} isAdminApp={true} />} />
                            <Route path={signupPath} render={(props) => <SignupPage {...props} isAdminApp={true} />}  />
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

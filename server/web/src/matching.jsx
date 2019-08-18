import React from 'react';
import { Button, ButtonToolbar, Dropdown, DropdownButton, Container, Table} from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';
import {matchRoundApiModule, deleteMatchRoundApiModule} from './api/match_round_api_module';
import {fetchGroupsApiModule} from './api/fetch_groups';
import {fetchMatchingRoundsApiModule} from './api/fetch_matching_rounds';
import { MODAL_TYPES, showAction} from './modal_container';
import {getCurrentGroup} from './group_context_reducer';
import GroupSelector from './group_selector';

const FETCH_MATCHING_ROUNDS_FOR_GROUP = "FETCH_MATCHING_ROUNDS_FOR_GROUP";
const FETCHING_MATCHING_ROUNDS_FOR_GROUP = "FETCHING_MATCHING_ROUNDS_FOR_GROUP";
const FETCHED_MATCHING_ROUNDS_FOR_GROUP = "FETCHED_MATCHING_ROUNDS_FOR_GROUP";
const ERROR_FETCHING_MATCHING_ROUNDS_FOR_GROUP = "ERROR_FETCHING_MATCHING_ROUNDS_FOR_GROUP";

const initialState = {
}

export function shouldFetchMatchingRoundsForGroupAction(group) {
    return {type: FETCH_MATCHING_ROUNDS_FOR_GROUP, groupToFetch: group};
}

export function fetchingMatchingRoundsForGroupAction() {
    return {type: FETCHING_MATCHING_ROUNDS_FOR_GROUP};
}

export function fetchedMatchingRoundsForGroupAction(matchingRounds) {
    return {type: FETCHED_MATCHING_ROUNDS_FOR_GROUP, matchingRounds: matchingRounds};
}

export function errorFetchingMatchingRoundsForGroupAction(error) {
    return {type: ERROR_FETCHING_MATCHING_ROUNDS_FOR_GROUP, fetchingMatchingRoundsForGroupError: error};
}

export function getShouldFetchMatchingRoundsForGroup(state) {
    return state.matchingReducer.shouldFetchMatchingRounds
}

export function getMatchingRounds(state) {
    return state.matchingReducer.matchingRounds;  
}

export function matchingReducer(state = initialState, action) {
    switch (action.type) {
        case FETCH_MATCHING_ROUNDS_FOR_GROUP:
            return Object.assign({}, state, {shouldFetchMatchingRounds: true, groupToFetch: action.groupToFetch});
        case FETCHING_MATCHING_ROUNDS_FOR_GROUP:
            return Object.assign({}, state, {shouldFetchMatchingRounds: false});
        case FETCHED_MATCHING_ROUNDS_FOR_GROUP:
            return Object.assign({}, state, {shouldFetchMatchingRounds: false, matchingRounds: action.matchingRounds});
        case ERROR_FETCHING_MATCHING_ROUNDS_FOR_GROUP:
            return Object.assign({}, state, {shouldFetchMatchingRounds: false, matchingRoundFetchError: action.fetchingMatchingRoundsForGroupError});
        default:
            return state;
    };
}
/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
 */
export class MatchingPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    componentDidMount() {
        if (!!this.props.groupToFetch) {
            this.props.fetchMatchingRoundsForGroup(this.props.groupToFetch);
        }
    }

    render() {
        return (
            <Container className="panel-body">
                <GroupSelector /> 
                <div className="panel-content">
                    <ButtonToolbar>
                        <Button variant="primary" size="lg" onClick={() => this.props.showModal(MODAL_TYPES.CREATE_MATCHING_ROUND)}>New Matching Round</Button>
                        {/* <Button variant="primary" size="lg">Delete members</Button>
                        <Button variant="primary" size="lg">Notify members</Button> */}
                    </ButtonToolbar>
                    <h2 className="mt-3">Matching Rounds</h2>
                    <div className="main-table-container">
                        <Table striped bordered hover>
                            <thead>
                                <tr>
                                    <th>#</th>
                                    <th>Name</th>
                                    <th>Stage</th>
                                </tr>
                            </thead>
                            <tbody>
                                {this.props.matchingRounds.map((matchingRound => {
                                    return (
                                        <tr>
                                            <td>{matchingRound.matchRoundId}</td>
                                            <td>{matchingRound.name}</td>
                                            <td>{matchingRound.status}</td>
                                        </tr>
                                    );
                                }))}
                                
                            </tbody>
                        </Table>
                    </div>
                </div>
            </Container>
        );
    }
}

const MatchingPageComponent = apiServiceConnect(
    (state) => ({
        groupToFetch: getCurrentGroup(state),
        groups: fetchGroupsApiModule.isFinished(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [],
        matchingRounds:  fetchMatchingRoundsApiModule.isFinished(state) ? fetchMatchingRoundsApiModule.getData(state): undefined || [],
    }),
    (dispatch) => ({
        fetchMatchingRoundsForGroup: (group) => dispatch(fetchMatchingRoundsApiModule.getApiExecuteAction({groupId: group.groupId})),
        deleteMatchingRound: (matchRoundId) => dispatch(deleteMatchRoundApiModule.getApiExecuteAction({matchRoundId})),
        showModal: (modalType) => dispatch(showAction(modalType))
    })
)(CookieAwareComponent(withCookies(MatchingPage)));

export default MatchingPageComponent;

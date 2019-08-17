import React from 'react';
import { Button, ButtonToolbar, Dropdown, DropdownButton, Container, Table} from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';
import {matchRoundApiModule, deleteMatchRoundApiModule} from './api/match_round_api_module';
import {fetchGroupsAction, getGroupsForAdmin} from './get_managed_groups_view';
import { MODAL_TYPES, showAction} from './modal_container';

const FETCH_MATCHING_ROUNDS_FOR_GROUP = "FETCH_MATCHING_ROUNDS_FOR_GROUP";
const FETCHING_MATCHING_ROUNDS_FOR_GROUP = "FETCHING_MATCHING_ROUNDS_FOR_GROUP";
const FETCHED_MATCHING_ROUNDS_FOR_GROUP = "FETCHED_MATCHING_ROUNDS_FOR_GROUP";
const ERROR_FETCHING_MATCHING_ROUNDS_FOR_GROUP = "ERROR_FETCHING_MATCHING_ROUNDS_FOR_GROUP";

const initialState = {
    shouldFetchMatchingRounds: false,
    groupToFetch: undefined,
    fetchingMatchingRoundsForGroupError: undefined,
    matchingRounds: undefined,
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

export function getMatchingRoundsGroupToFetch(state) {
    return state.matchingReducer.groupToFetch;
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
        this.onDropdownChanged = this.onDropdownChanged.bind(this);
    }

    componentDidMount() {
        this.props.fetchGroups();
        // TODO(skong, acod): use this as template
        this.props.createNewMatchingRoundForGroup("749625d4-a6b9-11e9-9737-0242ac130002", [1, 3], 1, 1, 2019);
        // this.props.deleteMatchingRound(1);
    }

    onDropdownChanged(group) {
        console.log("[onDropdownChanged] Change to group " + group.groupId)
        this.props.fetchMatchingRoundsForGroup(group);
    }

    render() {
        const dropdownItems = this.props.groups ? this.props.groups.map(group => <Dropdown.Item onClick={() => this.onDropdownChanged(group)} key={group.groupId} eventKey={group.groupId}> {group.groupName} </Dropdown.Item>) : [];
        console.log(this.props.matchingRounds)
        return (
            <Container className="panel-body">
                <div className="group-info">
                    <h2>You are currently managing: </h2>
                    <ButtonToolbar>
                        <DropdownButton
                            title={this.props.groupToFetch ? this.props.groupToFetch.groupName : 'Your Groups'}
                            variant='Primary'
                            id='managed-groups-dropdown'
                        >
                            {dropdownItems}
                        </DropdownButton>
                    </ButtonToolbar>
                </div>
                <div className="panel-content">
                    <ButtonToolbar>
                        <Button variant="primary" size="lg" onClick={() => this.props.showModal(MODAL_TYPES.CREATE_MATCHING_ROUND)}>New Matching Round</Button>
                        {/* <Button variant="primary" size="lg">Delete members</Button>
                        <Button variant="primary" size="lg">Notify members</Button> */}
                    </ButtonToolbar>
                    <h2 className="mt-3">Matching Rounds</h2>
                    <div className="members-table-container">
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
        groupToFetch: getMatchingRoundsGroupToFetch(state),
        groups: getGroupsForAdmin(state),
        matchingRounds: getMatchingRounds(state) || [],
    }),
    (dispatch) => ({
        fetchMatchingRoundsForGroup: (group) => dispatch(shouldFetchMatchingRoundsForGroupAction(group)),
        fetchGroups: () => dispatch(fetchGroupsAction()),
        // TODO(skong): use this
        createNewMatchingRoundForGroup: (groupId, userIds, maxLowerYearsPerUpperYear, maxUpperYearsPerLowerYear, youngestUpperGradYear) => dispatch(matchRoundApiModule.getApiExecuteAction({groupId, userIds, maxLowerYearsPerUpperYear, maxUpperYearsPerLowerYear, youngestUpperGradYear})),
        deleteMatchingRound: (matchRoundId) => dispatch(deleteMatchRoundApiModule.getApiExecuteAction({matchRoundId})),
        showModal: (modalType) => dispatch(showAction(modalType))
    })
)(CookieAwareComponent(withCookies(MatchingPage)));

export default MatchingPageComponent;

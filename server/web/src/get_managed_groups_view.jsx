import React from 'react';
import {connect, dispatch} from 'react-redux';
import {Accordion, Card, Button} from 'react-bootstrap';

import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';

const GOT_GROUPS = 'GOT_GROUPS';
const FETCH_GROUPS = 'FETCH_GROUPS';
const FETCHING_GROUP = 'FETCHING_GROUPS';
const ERROR_FETCHING_GROUPS = 'ERROR_FETCHING_GROUPS';

const initialState = {
    shouldFetchGroups: false,
    fetchingGroups: false,
    errorMessage: undefined,
    groups: []
}

export function getShouldFetchGroups(state) {
    return state.getManagedGroupsReducer.shouldFetchGroups;
}

export function gotGroupsAction(groups) {
    return {type: GOT_GROUPS, groups: groups};
}

export function fetchingGroupsAction() {
    return {type: FETCHING_GROUP};
}

export function fetchGroupsAction() {
    return {type: FETCH_GROUPS};
}

export function errorFetchingGroupsAction(errorMessage) {
    return {type: ERROR_FETCHING_GROUPS, errorMessage: errorMessage};
}

export function getManagedGroupsReducer(state = initialState, action) {
    switch(action.type) {
        case FETCH_GROUPS:
            return Object.assign({}, state, {shouldFetchGroups: true}); 
        case FETCHING_GROUP:
            return Object.assign({}, state, {fetchingGroups: true, shouldFetchGroups: false}); 
        case GOT_GROUPS:
            return Object.assign({}, state, {fetchingGroups: false, groups: action.groups});
        default:
            return state;
    }
}

class GetManagedGroupsView extends React.Component {
    componentDidMount() {
        this.props.fetchGroups();
    }

    render() {
        return (
            <div>
                <h2>Groups you manage</h2>
                    {this.props.groups.map((group) => {
                        return (
                        <Accordion key={group.managedGroupReferralEmail}>
                            <Card>
                                <Card.Header>
                                    <Accordion.Toggle as={Button} variant="link" eventKey="0">
                                        {group.groupName}
                                    </Accordion.Toggle>
                                </Card.Header>
                                <Accordion.Collapse eventKey="0">
                                    <Card.Body>{group.managedGroupReferralEmail}</Card.Body>
                                </Accordion.Collapse>
                            </Card>
                        </Accordion>
                        )
                    })}
            <p>{this.props.errorMessage}</p>
            </div >
        );
    }
}

const GetManagedGroupsViewComponent = connect(
    (state) => ({groups: state.getManagedGroupsReducer.groups || [], errorMessage: state.getManagedGroupsReducer.errorMessage}),
    (dispatch) => {
        return {
            gotGroups: (groups) => dispatch(gotGroupsAction(groups)),
            fetchGroups: () => dispatch(fetchGroupsAction())
        }
    }
)(CookieAwareComponent(withCookies(GetManagedGroupsView)));

export default GetManagedGroupsViewComponent;
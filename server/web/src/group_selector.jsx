import ReactDOM from 'react-dom';
import { Redirect } from 'react-router-dom';
import { ButtonToolbar, DropdownButton } from 'react-bootstrap';
import React from 'react';
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';

const GOT_GROUPS = 'GOT_GROUPS';
const FETCH_GROUPS = 'FETCH_GROUPS';
const FETCHING_GROUPS = 'FETCHING_GROUPS';
const ERROR_FETCHING_GROUPS = 'ERROR_FETCHING_GROUPS';

const initialState = {
    selectedGroupIndex: 0,
    groups: [],
    shouldFetchGroups: false,
    fetchingGroups: false,
    errorMessage: undefined
}

export function fetchingGroupsAction() {
    return {type: FETCHING_GROUPS};
}

export function fetchGroupsAction() {
    return {type: FETCH_GROUPS};
}

export function gotGroupsAction(groups) {
    console.log(groups);
    return {type: GOT_GROUPS, groups: groups};
}

export function errorFetchingGroupsAction(errorMessage) {
    return {type: ERROR_FETCHING_GROUPS, errorMessage: errorMessage};
}

export function groupsReducer(state = initialState, action) {
    switch(action.type) {
        case FETCH_GROUPS:
            return Object.assign({}, state, {shouldFetchGroups: true}); 
        case FETCHING_GROUPS:
            return Object.assign({}, state, {shouldFetchGroups: false, fetchingGroups: true}); 
        case GOT_GROUPS:
            return Object.assign({}, state, {shouldFetchGroups: false, groups: action.groups});
        default:
            return state;
    }
}

class GroupSelector extends React.Component {
    // populateMe() {
    //     if (!!this.props.isAuthenticated) {
    //         HiveApiService.me(
    //             ({ Result }) => {
    //                 if (!this.state.me || this.state.me.userId != Result.userId ) {
    //                     this.setState({ me: {
    //                         userId: Result.userId,
    //                         firstName: Result.firstName,
    //                         lastName: Result.lastName,
    //                         email: Result.email
    //                     }});
    //                 }
    //             },
    //             err => console.log
    //         );
    //     }
    // }

    // componentDidMount() {
    //     this.populateMe();
    // }

    // componentDidUpdate() {
    //     this.populateMe();
    // }

    render() {
        const dropdownItems = this.props.groupState.groups.map(group => <Dropdown.Item onClick={() => this.onDropdownChanged(group)} key={group.groupId} eventKey={group.groupId}> {group.groupName} </Dropdown.Item>)
        return (
          <div className="group-info">
            <h2>You are currently managing: </h2>
            <ButtonToolbar>
                <DropdownButton
                    title='Your Groups'
                    variant='Primary'
                    id='managed-groups-dropdown'
                >
                    {dropdownItems}
                </DropdownButton>
            </ButtonToolbar>
          </div>
        );
    }
}

const GroupSelectorComponent = connect(
    (state) => {
        return {
            groupState: state.groupsReducer
        };
    },
    null
    // (dispatch) => {
    //     return {
    //         didLogout: (state) => {dispatch(logoutAction(state))}
    //     };
    // }
)(CookieAwareComponent(withCookies(GroupSelector)));

export default GroupSelectorComponent;
import React from 'react';
import { Container} from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx';
import { showAction } from './modal_container';
import {withCookies} from 'react-cookie';
import { HiveApiService } from './api_controller.js';
import apiServiceConnect from './api/api_service_connect';
import { gotGroupsAction, fetchGroupsAction } from './get_managed_groups_view';

// const GROUPS = ['Hello Kitty', 'My Little Unicorn', 'Black Mamba'];
const STATS = ['200 members', '20 unregistered', '180 registered', '180 matched'];

const GOT_MEMBERS = 'GOT_MEMBERS';
const FETCH_MEMBERS = 'FETCH_MEMBERS';
const FETCHING_MEMBERS = 'FETCHING_MEMBERS';
const ERROR_FETCHING_MEMBERS= 'ERROR_FETCHING_MEMBERS';

const initialState = {
    shouldFetchMembers: false,
    fetchingMembers: false,
    errorMessage: undefined,
    members: []
}

export function getShouldFetchMembers(state) {
    return state.membersReducer.shouldFetchMembers;
}

export function gotMembersAction(members) {
    return {type: GOT_MEMBERS, members: members};
}

export function fetchingMembersAction() {
    return {type: FETCHING_MEMBERS};
}

export function fetchMembersAction() {
    return {type: FETCH_MEMBERS};
}

export function errorFetchingMembersAction(errorMessage) {
    return {type: ERROR_FETCHING_MEMBERS, errorMessage: errorMessage};
}

export function membersReducer(state = initialState, action) {
    switch(action.type) {
        case FETCH_MEMBERS:
            return Object.assign({}, state, {shouldFetchMembers: true}); 
        case FETCHING_MEMBERS:
            return Object.assign({}, state, {fetchingMembers: true, shouldFetchMembers: false}); 
        case GOT_MEMBERS:
            return Object.assign({}, state, {fetchingMembers: false, members: action.members});
        default:
            return state;
    }
}

const GROUPS = ['Hello Kitty', 'My Little Unicorn', 'Black Mamba'];
const STATS = ['200 members', '20 unregistered', '180 registered', '180 matched'];

/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
 */

export class MembersPage extends React.Component {

    componentDidMount() {
        this.props.fetchMembers();
    }

    // populateGroups() {
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

    render() {
        // const dropdownItems = GROUPS.map((groupName, i) => <Dropdown.Item key={i} eventKey={i.toString()}> {groupName} </Dropdown.Item>)
        console.log(this.props.groups);
        console.log(this.props.members);
        const dropdownItems = this.props.groups.map((group, i) => <Dropdown.Item key={i} eventKey={i.toString()}> {group.groupName} </Dropdown.Item>)
        const statItems = STATS.map((stat, i) => <div key={i} className="members-stat"> {stat} </div>)
        return (
            <Container className="panel-body">
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
                <div className="panel-content">
                    <ButtonToolbar>
                        <Button variant="primary" size="lg" onClick={this.props.showModal}>Add members</Button>
                        <Button variant="primary" size="lg" onClick={this.props.showModal}>Delete members</Button>
                        <Button variant="primary" size="lg" onClick={this.props.showModal}>Notify members</Button>
                    </ButtonToolbar>
                    <h2 className="mt-3">Members</h2>
                    <div className="members-stats-container">
                        {statItems}
                    </div>
                    <div className="members-table-container">
                        <Table striped bordered hover>
                            <thead>
                                <tr>
                                    <th>#</th>
                                    <th>Name</th>
                                    <th>Status</th>
                                    <th>Program</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td>1</td>
                                    <td>Mark</td>
                                    <td>Otto</td>
                                    <td>@mdo</td>
                                </tr>
                                <tr>
                                    <td>2</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                            </tbody>
                        </Table>
                    </div>
                </div>
            </Container>
        );
    }
}

const MembersPageComponent = apiServiceConnect(
    (state) => ({
        groups: state.getManagedGroupsReducer.groups || [{groupName: "Cannot fetch groups"}], 
        members: state.membersReducer.members || [],
        errorMessage: state.getManagedGroupsReducer.errorMessage,
        // TODO: rename
        membersErrorMessage: state.membersReducer.errorMessage
    }),
    (dispatch) => {
        return {
            gotGroups: (groups) => dispatch(gotGroupsAction(groups)),
            fetchGroups: () => dispatch(fetchGroupsAction()),
            showModal: (state) => dispatch(showAction(state)),
            gotMembers: (members) => dispatch(gotMembersAction(members)),
            fetchMembers: () => dispatch(fetchMembersAction())
        }
    }
)(CookieAwareComponent(withCookies(MembersPage)));

export default MembersPageComponent;

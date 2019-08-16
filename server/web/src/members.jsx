import React from 'react';
import { Container, ButtonToolbar, Dropdown, DropdownButton, Button, Table } from "react-bootstrap";
import {BootstrapTable} from 'react-bootstrap-table';
import CookieAwareComponent from './cookie_aware_component.jsx';
import { MODAL_TYPES, showAction} from './modal_container';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';
import { gotGroupsAction, fetchGroupsAction, getGroupsForAdmin } from './get_managed_groups_view';
import {userGroupDeleteApiModule} from './api/user_group_delete_api_module'

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

export function getMembersFromState(state) {
    return state.membersReducer.members;
}

export function gotMembersAction(members) {
    return {type: GOT_MEMBERS, members: members};
}

export function getGroupToFetch(state) {
    return state.membersReducer.groupToFetch;
}

export function fetchingMembersAction() {
    return {type: FETCHING_MEMBERS};
}

export function fetchMembersAction(group) {
    return {type: FETCH_MEMBERS, group: group};
}

export function errorFetchingMembersAction(errorMessage) {
    return {type: ERROR_FETCHING_MEMBERS, errorMessage: errorMessage};
}

export function membersReducer(state = initialState, action) {
    switch(action.type) {
        case FETCH_MEMBERS:
            return Object.assign({}, state, {shouldFetchMembers: true, groupToFetch: action.group}); 
        case FETCHING_MEMBERS:
            return Object.assign({}, state, {shouldFetchMembers: false, fetchingMembers: true}); 
        case GOT_MEMBERS:
            return Object.assign({}, state, {shouldFetchMembers: false, fetchingMembers: false, members: action.members});
        default:
            return state;
    }
}

/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
 */

export class MembersPage extends React.Component {

    constructor(props) {
        super(props);
        this.state =  {
            selected: [],
        };
        this.onDropdownChanged = this.onDropdownChanged.bind(this);
        this.onRowSelect = this.onRowSelect.bind(this);
        this.selectRowProp = {
            mode: 'checkbox',
            onSelect: this.onRowSelect,
        }
    }

    componentDidMount() {
        // kickoff initial fetch
        // this.props.fetchGroups();
        this.props.fetchMembers();
        // TODO(skong, acod): use this as a template
        // this.props.deleteMemberFromGroup(1, "9ba4177a-a6b8-11e9-81f1-0242ac130002");
    }

    onDropdownChanged(group) {
        console.log("[onDropdownChanged] Change to group " + group.groupId)
        this.props.fetchMembers(group);
    }

    onRowSelect({ id }, isSelected) {
        if (isSelected) {
          this.setState({
            selected: [ ...this.state.selected, id ],
          });
        } else {
           this.setState({
               selected: [ ...this.state.selected.filter(idToCompare => idToCompare != id)],
           }); 
        }
        return ;
    }

    render() {
        const dropdownItems = this.props.groups.map(group => <Dropdown.Item onClick={() => this.onDropdownChanged(group)} key={group.groupId} eventKey={group.groupId}> {group.groupName} </Dropdown.Item>)
        const statItems = STATS.map((stat, i) => <div key={i} className="members-stat"> {stat} </div>)
        return (
            <Container className="panel-body">
                <div className="group-info">
                    <h2>You are currently managing: </h2>
                    <ButtonToolbar>
                        <DropdownButton
                            title={ this.props.groupToFetch ? this.props.groupToFetch.groupName : undefined || 'Your Groups'}
                            variant='Primary'
                            id='managed-groups-dropdown'
                        >
                            {dropdownItems}
                        </DropdownButton>
                    </ButtonToolbar>
                </div>
                <div className="panel-content">
                    <ButtonToolbar>
                        <Button variant="primary" size="lg" onClick={() => this.props.showModal(MODAL_TYPES.ADD_MEMBER)}>Add members</Button>
                        <Button variant="primary" size="lg" onClick={() => this.props.showModal(MODAL_TYPES.DELETE_MEMBER)}>Delete members</Button>
                        {/* <Button variant="primary" size="lg" onClick={this.props.showModal}>Notify members</Button> */}
                    </ButtonToolbar>
                    <h2 className="mt-3">Members</h2>
                    <div className="members-stats-container">
                        {statItems}
                    </div>
                    <div className="members-table-container">
                        <BootstrapTable data={this.props.members.map(groupMember => {
                                    return ({
                                        id: groupMember.user.userId,
                                        name: groupMember.user.firstName + " " + groupMember.user.lastName,
                                        status: groupMember.status,
                                        email: groupMember.email,
                                        programName: groupMember.cohort ? (groupMember.cohort.programName + " " + groupMember.cohort.gradYear) : "No cohort"
                                    });
                                })} selectRow={this.selectRowProp}>
                            <TableHeaderColumn dataField='id' isKey>User Id</TableHeaderColumn>
                            <TableHeaderColumn dataField='name'>User Name</TableHeaderColumn>
                            <TableHeaderColumn dataField='email'>User Email</TableHeaderColumn>
                            <TableHeaderColumn dataField='status'>Status</TableHeaderColumn>
                            <TableHeaderColumn dataField='programName'>Program</TableHeaderColumn>
                        </BootstrapTable>
                    </div>
                </div>
            </Container>
        );
    }
}

const MembersPageComponent = apiServiceConnect(
    (state) => ({
        groupToFetch: getGroupToFetch(state),
        groups: getGroupsForAdmin(state) || [], 
        members: getMembersFromState(state) || [],
        errorMessage: state.getManagedGroupsReducer.errorMessage,
        // TODO: rename
        membersErrorMessage: state.membersReducer.errorMessage
    }),
    (dispatch) => {
        return {
            gotGroups: (groups) => dispatch(gotGroupsAction(groups)),
            fetchGroups: () => dispatch(fetchGroupsAction()),
            showModal: (modalType) => dispatch(showAction(modalType)),
            gotMembers: (members) => dispatch(gotMembersAction(members)),
            fetchMembers: (groupId) => dispatch(fetchMembersAction(groupId)),
            // TODO(skong): use this
            deleteMemberFromGroup: (userId, groupId) => dispatch(userGroupDeleteApiModule.getApiExecuteAction({userId, groupId})),
        }
    }
)(CookieAwareComponent(withCookies(MembersPage)));

export default MembersPageComponent;

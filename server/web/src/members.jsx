import React from 'react';
import { Container, ButtonToolbar, Dropdown, DropdownButton, Button, Table } from "react-bootstrap";
import {BootstrapTable, TableHeaderColumn} from 'react-bootstrap-table';
import CookieAwareComponent from './cookie_aware_component.jsx';
import { MODAL_TYPES, showAction} from './modal_container';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';
import { fetchGroupsApiModule } from './api/fetch_groups';
import {userGroupDeleteApiModule} from './api/user_group_delete_api_module';
import GroupSelector from './group_selector';

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
            shouldRefresh: true,
        };
        this.onDropdownChanged = this.onDropdownChanged.bind(this);
        this.onRowSelect = this.onRowSelect.bind(this);
        this.onSelectAll = this.onSelectAll.bind(this);
        this.selectRowProp = {
            mode: 'checkbox',
            onSelect: this.onRowSelect,
            onSelectAll: this.onSelectAll,
        }
        this.deleteSelectedUsers = this.deleteSelectedUsers.bind(this);
    }

    componentDidMount() {
        // kickoff initial fetch
        this.props.fetchGroups();
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

    onSelectAll(isSelected, rows) {
        if (isSelected) {
            this.setState({
                selected: [...this.state.selected, rows],
            });
        } else {
            this.setState({
                selected: [],
            });
        }
    }

    deleteSelectedUsers() {
        this.state.selected.forEach(userId => this.props.deleteMemberFromGroup(userId, this.props.groupToFetch.groupId));
        this.setState({
            shouldRefresh: true,
        });
    }

    render() {
        if (this.state.shouldRefresh) {
            this.props.fetchMembers();
            this.setState({shouldRefresh: false});
        }
        const dropdownItems = this.props.groups.map(group => <Dropdown.Item onClick={() => this.onDropdownChanged(group)} key={group.groupId} eventKey={group.groupId}> {group.groupName} </Dropdown.Item>)
        const statItems = STATS.map((stat, i) => <div key={i} className="members-stat"> {stat} </div>)
        return (
            <Container className="panel-body">
                <GroupSelector/>
                <div className="panel-content">
                    <ButtonToolbar>
                        <Button variant="primary" size="lg" onClick={() => this.props.showModal(MODAL_TYPES.ADD_MEMBER)}>Add members</Button>
                        <Button variant="primary" size="lg" onClick={this.deleteSelectedUsers}>Delete members</Button>
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
        groups: fetchGroupsApiModule.isFinished(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [], 
        members: getMembersFromState(state) || [],
        errorMessage: fetchGroupsApiModule.getErrorMessage(state),
        didCompleteDelete: userGroupDeleteApiModule.isFinished(state),
        // TODO: rename
        membersErrorMessage: state.membersReducer.errorMessage
    }),
    (dispatch) => {
        return {
            fetchGroups: () => dispatch(fetchGroupsApiModule.getApiExecuteAction()),
            showModal: (modalType) => dispatch(showAction(modalType)),
            gotMembers: (members) => dispatch(gotMembersAction(members)),
            fetchMembers: (groupId) => dispatch(fetchMembersAction(groupId)),
            deleteMemberFromGroup: (userId, groupId) => dispatch(userGroupDeleteApiModule.getApiExecuteAction({userId, groupId})),
        }
    }
)(CookieAwareComponent(withCookies(MembersPage)));

export default MembersPageComponent;

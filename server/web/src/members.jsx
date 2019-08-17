import React from 'react';
import { Container, ButtonToolbar, Dropdown, DropdownButton, Button, Table } from "react-bootstrap";
import BootstrapTable from 'react-bootstrap-table-next';
import CookieAwareComponent from './cookie_aware_component.jsx';
import { MODAL_TYPES, showAction} from './modal_container';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';
import {fetchGroupsApiModule} from './api/fetch_groups';
import {fetchMembersApiModule} from './api/fetch_members';
import {userGroupDeleteApiModule} from './api/user_group_delete_api_module'
import {getCurrentGroup} from './group_context_reducer';
import GroupSelector from './group_selector';

// const GROUPS = ['Hello Kitty', 'My Little Unicorn', 'Black Mamba'];
const STATS = ['200 members', '20 unregistered', '180 registered', '180 matched'];

/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
 */
export class MembersPage extends React.Component {

    constructor(props) {
        super(props);
        this.state =  {
            selected: [],
            shouldRefresh: false,
        };
        this.onRowSelect = this.onRowSelect.bind(this);
        this.onSelectAll = this.onSelectAll.bind(this);
        this.selectRowProp = {
            mode: 'checkbox',
            clickToSelect: true,
            onSelect: this.onRowSelect,
            onSelectAll: this.onSelectAll,
        }
        this.deleteSelectedUsers = this.deleteSelectedUsers.bind(this);
    }

    componentDidMount() {
        if (!!this.props.groupToFetch) {
            this.props.fetchMembers(this.props.groupToFetch.groupId);
        }
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
            this.props.fetchMembers(this.props.groupToFetch.groupId);
            this.setState({shouldRefresh: false});
        }
        const statItems = STATS.map((stat, i) => <div key={i} className="members-stat"> {stat} </div>)

        const membersData = this.props.members ?  this.props.members.map(groupMember => {
            return ({
                id: groupMember.user.userId,
                name: groupMember.user.firstName + " " + groupMember.user.lastName,
                status: groupMember.status,
                email: groupMember.email,
                programName: groupMember.cohort ? (groupMember.cohort.programName + " " + groupMember.cohort.gradYear) : "No cohort"
            });
        }): [];
        const columns = [
            {
                dataField: 'id',
                text: 'User ID'
            },
            {
                dataField: 'name',
                text: 'Name'
            },
            {
                dataField: 'email',
                text: 'Email'
            },
            {
                dataField: 'status',
                text: 'Status'
            },
            {
                dataField: 'programName',
                text: 'Program'
            }
        ];
        return (
            <Container className="panel-body">
                <GroupSelector /> 
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
                        <BootstrapTable keyField='id' data={membersData} columns={columns} selectRow={this.selectRowProp}>
                        </BootstrapTable>
                    </div>
                </div>
            </Container>
        );
    }
}

const MembersPageComponent = apiServiceConnect(
    (state) => ({
        groupToFetch: getCurrentGroup(state),
        groups: fetchGroupsApiModule.isFinished(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [], 
        members: fetchMembersApiModule.isFinished(state) ? fetchMembersApiModule.getData(state): undefined || [],
        errorMessage: fetchGroupsApiModule.getErrorMessage(state),
        didCompleteDelete: userGroupDeleteApiModule.isFinished(state),
        // TODO: rename
        membersErrorMessage: fetchMembersApiModule.getErrorMessage(state),
    }),
    (dispatch) => {
        return {
            showModal: (modalType) => dispatch(showAction(modalType)),
            fetchMembers: (groupId) => dispatch(fetchMembersApiModule.getApiExecuteAction({groupId})),
            deleteMemberFromGroup: (userId, groupId) => dispatch(userGroupDeleteApiModule.getApiExecuteAction({userId, groupId})),
        }
    }
)(CookieAwareComponent(withCookies(MembersPage)));

export default MembersPageComponent;

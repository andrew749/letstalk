import React from 'react';
import { ButtonToolbar, Dropdown, DropdownButton} from "react-bootstrap";
import apiServiceConnect from './api/api_service_connect';
import {fetchGroupsApiModule} from './api/fetch_groups';
import {setNewCurrentGroupAction, getCurrentGroup} from './group_context_reducer';

/**
 * Props:
 *  listeners: callbacks to execute on change passed the new group
 */
class GroupSelector extends React.Component {
    constructor(props) {
        super(props);
        this.onDropdownChanged = this.onDropdownChanged.bind(this);
    }

    componentDidMount() {
        this.props.fetchGroups();
    }

    onDropdownChanged(group) {
        console.log("[onDropdownChanged] Change to group " + group.groupId)
        this.props.updateCurrentGroup(group);
        this.props.listeners.forEach(listener => listener(group));
    }

    render() {
        let dropdownItems = [];
        if (this.props.groups) {
            dropdownItems = this.props.groups.map(group => <Dropdown.Item onClick={() => this.onDropdownChanged(group)} key={group.groupId} eventKey={group.groupId}> {group.groupName} </Dropdown.Item>)
        }
        return (
            <div className="group-info">
                <h2>You are currently managing: </h2>
                <ButtonToolbar>
                    <DropdownButton
                        title={(this.props.groupToFetch ? this.props.groupToFetch.groupName : undefined) || 'Your Groups'}
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

export default apiServiceConnect(
    (state) => ({
        groupToFetch: getCurrentGroup(state),
        groups: fetchGroupsApiModule.getData(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [],
    }),
    (dispatch) => ({
        fetchGroups: () => dispatch(fetchGroupsApiModule.getApiExecuteAction()),
        updateCurrentGroup: (group) => dispatch(setNewCurrentGroupAction(group))
    }) 
)(GroupSelector);
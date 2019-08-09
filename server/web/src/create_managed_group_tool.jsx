import React from 'react';
import {Form, Button} from 'react-bootstrap';

import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';

import {onChange} from './util.js';
import { fetchGroupsAction } from './get_managed_groups_view.jsx';
import apiServiceConnect from './api/api_service_connect';

const CREATE_GROUP = 'CREATE_GROUP';

const initialState = {
    creatingGroups: false
}

export function createGroupAction() {
    return {type: CREATE_GROUP};
}

export function groupManagementReducer(state = initialState, action) {
    switch(action.type) {
        case CREATE_GROUP:
            return Object.assign({}, state, {creatingGroups: true});
        default:
            return state;
    }
}

class CreateManagedGroupTool extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            createManagedGroupToolModel: {
                groupName: "",
            }
        }

        this.onChange = onChange.bind(this, 'createManagedGroupToolModel');
        this.createGroup = this.createGroup.bind(this);
    }

    createGroup(e) {
        e.preventDefault();
        e.stopPropagation();
        const model = this.state.createManagedGroupToolModel;
        this.props.apiService.createManagedGroup(model.groupName, () => {/* When we start creating */}, () => {this.props.refreshPage()}, (err) => {console.warn(err);});
    }

    render() {
        return (
            <div>
                <h2>Create New Managed Group</h2>
                <Form onSubmit={this.createGroup}>
                    <Form.Group controlId="formCreateGroup">
                        <Form.Control size="lg" type="text" name="groupName" placeholder="Group Name" onChange={this.onChange} />
                    </Form.Group>
                    <Button variant="primary" type="submit">
                        Create Group
                    </Button>
                </Form>
            </div>
        );
    }
}

const CreateManagedGroupToolComponent = apiServiceConnect(
    null,
    (dispatch) => {
        return {
            createGroup: () => dispatch(createGroupAction()),
            refreshPage: () => dispatch(fetchGroupsAction())
        }
    }
)(CookieAwareComponent(withCookies(CreateManagedGroupTool)));

export default CreateManagedGroupToolComponent;
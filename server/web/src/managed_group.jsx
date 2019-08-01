import React from 'react';
import { Container } from 'react-bootstrap';

import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import CreateManagedGroupTool from './create_managed_group_tool.jsx';
import GetManagedGroupsView from './get_managed_groups_view.jsx';
import apiServiceConnect from './api/api_service_connect';

class ManagedGroupPage extends React.Component {
    render() {
        return (
            <Container className="panel-body">
                <CreateManagedGroupTool />
                <GetManagedGroupsView />
            </Container>
        );
    }
}

const ManagedGroupPageComponent = apiServiceConnect(
    null,
    null,
)(CookieAwareComponent(withCookies(ManagedGroupPage)));


export default ManagedGroupPageComponent;
import React from 'react';
import {connect} from 'react-redux';

import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import CreateManagedGroupTool from './create_managed_group_tool.jsx';
import GetManagedGroupsView from './get_managed_groups_view.jsx';

class ManagedGroupPage extends React.Component {
    render() {
        return (
            <div>
                <CreateManagedGroupTool />
                <GetManagedGroupsView />
            </div>
        );
    }
}

const ManagedGroupPageComponent = connect(
    null,
    null,
)(CookieAwareComponent(withCookies(ManagedGroupPage)));


export default ManagedGroupPageComponent;
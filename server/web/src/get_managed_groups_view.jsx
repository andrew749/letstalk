import React from 'react';
import {Accordion, Card, Button} from 'react-bootstrap';

import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';
import {fetchGroupsApiModule} from './api/fetch_groups';

class GetManagedGroupsView extends React.Component {
    componentDidMount() {
        this.props.fetchGroups();
    }

    render() {
        let groupsBody;
        if (!!this.props.groups && this.props.groups.length > 0) {
            groupsBody = (
                <div>
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
                </div>
            );
        } else {
            groupsBody = (
                <div>
                    <p>{"You currently manage no groups. Please create one first."}</p>
                </div>
            );
        }
        return (
            <div>
                {groupsBody}
                <p>{this.props.errorMessage}</p>
            </div >
        );
    }
}

const GetManagedGroupsViewComponent = apiServiceConnect(
    (state) => ({
        groups: fetchGroupsApiModule.getData(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [], 
        errorMessage: fetchGroupsApiModule.getErrorMessage(state),
    }),
    (dispatch) => {
        return {
            fetchGroups: () => dispatch(fetchGroupsApiModule.getApiExecuteAction())
        }
    }
)(CookieAwareComponent(withCookies(GetManagedGroupsView)));

export default GetManagedGroupsViewComponent;
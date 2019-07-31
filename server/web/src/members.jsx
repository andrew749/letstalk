import React from 'react';
import { Container} from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';

/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
 */
export class MembersPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {
        return (
            <Container>
                Members
            </Container>
        );
    }
}

const MembersPageComponent = apiServiceConnect()(CookieAwareComponent(withCookies(MembersPage)));

export default MembersPageComponent;

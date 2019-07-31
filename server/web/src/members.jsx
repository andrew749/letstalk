import React from 'react';
import { Redirect, Link} from 'react-router-dom';
import { Button, Container} from "react-bootstrap";
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';

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

const MembersPageComponent = connect()(CookieAwareComponent(withCookies(MembersPage)));

export default MembersPageComponent;

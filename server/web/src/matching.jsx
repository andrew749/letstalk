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
export class MatchingPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {
        return (
            <Container>
                Matching 
            </Container>
        );
    }
}

const MatchingPageComponent = connect()(CookieAwareComponent(withCookies(MatchingPage)));

export default MatchingPageComponent;

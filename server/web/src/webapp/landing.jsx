import React from 'react';
import { withCookies } from 'react-cookie';
import { connect } from 'react-redux';

import CookieAwareComponent from '../cookie_aware_component.jsx';


export class LandingPage extends React.Component {
    render() {
        return (
            <h1>Hive Landing page!</h1>
        );
    }
}

const LandingPageComponent = connect(
    null,
    (dispatch) => {
        return {};
    }
)(CookieAwareComponent(withCookies(LandingPage)));

export default LandingPageComponent;
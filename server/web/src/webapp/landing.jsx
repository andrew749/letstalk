import React from 'react';
import { withCookies } from 'react-cookie';
import { Container } from 'react-bootstrap';
import apiServiceConnect from '../api/api_service_connect';
import {bootstrapApiModule} from '../api/bootstrap_api_module';

import CookieAwareComponent from '../cookie_aware_component.jsx';

export class LandingPage extends React.Component {
    componentDidMount() {
        // this.props.getBootstrap();
    }

    render() {
        console.log(this.props.bootstrap);
        return (
            <Container>
                <h1>Hive Landing page!</h1>
            </Container>
        );
    }
}

const LandingPageComponent = apiServiceConnect(
    (state) => {
        return {
            bootstrap: bootstrapApiModule.getData(state)
        };
    },
    (dispatch) => {
        return {
            getBootstrap: () => dispatch(bootstrapApiModule.getApiExecuteAction({}))
        };
    }
)(CookieAwareComponent(withCookies(LandingPage)));

export default LandingPageComponent;

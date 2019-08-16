import React from 'react';
import { withCookies } from 'react-cookie';
import { Alert, Container } from 'react-bootstrap';
import { Redirect } from 'react-router-dom';
import ClipLoader from 'react-spinners/ClipLoader';

import apiServiceConnect from '../api/api_service_connect';
import {bootstrapApiModule} from '../api/bootstrap_api_module';
import {getFetchInfo} from '../api/api_module';
import {
    verifyEmailPathWeb,
    setCohortPathWeb,
    surveyPathWeb,
} from '../routes';

import CookieAwareComponent from '../cookie_aware_component.jsx';

export class LandingPage extends React.Component {
    componentDidMount() {
        this.props.getBootstrap();
    }

    renderSuccessBody() {
        switch (this.props.bootstrap.state) {
            case "account_created":
                return <Redirect to={verifyEmailPathWeb} />;
            case "account_email_verified":
                return <Redirect to={setCohortPathWeb} />;
            case "account_has_basic_info":
                return <Redirect to={surveyPathWeb} />;
            default:
                return <h3>"Hello"</h3>;
        }
    }

    render() {
        const { fetchState, errorMessage } = this.props.bootstrapFetchInfo;
        let body = null;
        if (fetchState === "success") {
            body = this.renderSuccessBody();
        } else if (fetchState === "error") {
            body = (
                <Alert variant="danger">
                    Failed to load your information with the following error: {errorMessage}
                </Alert>
            );
        } else {
            body = <ClipLoader />
        }
        return <Container>{ body }</Container>;
    }
}

const LandingPageComponent = apiServiceConnect(
    (state) => {
        return {
            bootstrap: bootstrapApiModule.getData(state),
            bootstrapFetchInfo: getFetchInfo(bootstrapApiModule, state),
        };
    },
    (dispatch) => {
        return {
            getBootstrap: () => dispatch(bootstrapApiModule.getApiExecuteAction({}))
        };
    }
)(CookieAwareComponent(withCookies(LandingPage)));

export default LandingPageComponent;

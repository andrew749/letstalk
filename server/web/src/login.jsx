import React from 'react';
import {Link} from 'react-router-dom';
import { Button, Container, FormGroup, FormControl, ControlLabel, Alert, Form } from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {getLandingPath, signupPathWeb, signupPath} from './routes.js';
import PostAuthRedirect from './post-auth-redirect.jsx';

import apiServiceConnect from './api/api_service_connect';
import HiveLogoWithTitle from './img/logo_with_title.png';

const LOGIN_ACTION = 'LOGIN';
const LOGOUT_ACTION = 'LOGOUT';

const initialState = {
  isAuthenticated: false,
  sessionId: undefined
}

export function loginAction(sessionId) {
    return {type: LOGIN_ACTION, sessionId: sessionId};
}

export function logoutAction() {
    return  {type: LOGOUT_ACTION};
}

export function loginReducer(state = initialState, action) {
    switch(action.type) {
        case LOGIN_ACTION:
            return Object.assign({}, state, {isAuthenticated: true, sessionId: action.sessionId});
        case LOGOUT_ACTION:
            return Object.assign({}, state, {isAuthenticated: false, sessionId: undefined});
        default:
            return state;
    }
}

/**
 * Props:
 *  - isAdminApp: determine whether this is the admin app
 */
export class LoginPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            redirectToReferrer: false,
            email: "",
            password: ""
        };

        this.handleChange = this.handleChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }

    validateForm() {
        return this.state.email.length > 0 && this.state.password.length > 0;
    }

    handleChange(event) {
        this.setState({
            [event.target.id]: event.target.value
        });
    }

    onSubmit(event) {
        event.preventDefault();
        // send to api server
        this.props.apiService.login(this.state.email, this.state.password)
            .then((data) => {
                this.props.didAuthenticate(data.Result.sessionId);
                this.setState({
                    submitState: 'SUCCESS',
                    redirectToReferrer: true
                });
            }).catch(err => {
                this.setState({
                    submitState: 'ERROR',
                    err: err.serverMessage
                });
            });
    }

    render() {
        const { redirectToReferrer } = this.state;
        const defaultPath = getLandingPath(this.props.isAdminApp);

        if (!!this.state.redirectToReferrer) {
            return <PostAuthRedirect defaultPath={defaultPath} />;
        }

        let alert;
        if (this.state.submitState) {
            if (this.state.submitState === "SUCCESS") {
                alert = (<Alert variant="success">Succesfully logged in.</Alert>);
            } else if (this.state.submitState === "ERROR") {
                alert = (<Alert variant="danger">Failed to login because {JSON.stringify(this.state.err)}</Alert>)
            }
        }

        let signupLink = null;
        if (!this.props.isAdminApp) {
            signupLink=(
                <div>
                    <Link to={(this.props.isAdminApp) ? signupPath : signupPathWeb}>
                        Don't have an account? Click to signup!
                    </Link>
                </div>
            );
        }
        return (
            <Container className="panel-body">
                <img id="hive-logo" src={HiveLogoWithTitle}/>
                <Form className="login-form" onSubmit={this.onSubmit}>
                    <Form.Group controlId="email">
                        <Form.Label>Email</Form.Label>
                        <Form.Control
                            autoFocus
                            type="email"
                            value={this.state.email}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="password">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            value={this.state.password}
                            onChange={this.handleChange}
                            type="password"
                        />
                    </Form.Group>
                    <Button
                        block
                        disabled={!this.validateForm()}
                        type="submit"
                    >
                        Login
                    </Button>
                </Form>
                <div className="message-container">
                    {alert}
                </div>
                {signupLink}
            </Container>
        );
    }
}

const LoginPageComponent = apiServiceConnect(
    null,
    (dispatch) => ({
        didAuthenticate: (state) => { dispatch(loginAction(state)) }
    }),
    undefined
)(CookieAwareComponent(withCookies(LoginPage)));

export default LoginPageComponent;

import React from 'react';
import {Redirect, Link} from 'react-router-dom';
import { Button, Container, FormGroup, FormControl, ControlLabel, Alert, Form } from "react-bootstrap";
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {landingPath, landingPathWeb, signupPathWeb, signupPath} from './routes.js';
import {HiveApiService} from './api_controller.js';

const LOGIN_ACTION = 'LOGIN';

const initialState = {
  isAuthenticated: false,
  sessionId: undefined
}

export function loginAction(sessionId) {
    return {type: LOGIN_ACTION, sessionId: sessionId};
}

export function loginReducer(state = initialState, action) {
    switch(action.type) {
        case LOGIN_ACTION:
            return Object.assign({}, state, {isAuthenticated: true, sessionId: action.sessionId});
        default:
            return state;
    }
}

/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
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
        const {cookies} = this.props;
        event.preventDefault();
        // send to api server

        HiveApiService.login(this.state.email, this.state.password)
            .then((data) => {
                cookies.set('sessionId', data.Result.sessionId);
                this.props.didAuthenticate(data.Result.sessionId);
                this.setState({
                    submitState: 'SUCCESS',
                    redirectToReferrer: true
                });
            }).catch(err => {
                this.setState({
                    submitState: 'ERROR',
                    err: err.body
                });
            });
    }

    render() {

        let { redirectToReferrer } = this.state;
        let { from } = this.props.location.state || { from: { pathname: this.props.isAdminPage ? landingPath : landingPathWeb } };

        if (!!redirectToReferrer) {
            return <Redirect to={from} />;
        }

        let alert;
        if (this.state.submitState) {
            if (this.state.submitState === "SUCCESS") {
                alert = (<Alert variant="success">Succesfully logged in.</Alert>);
            } else if (this.state.submitState === "ERROR") {
                alert = (<Alert variant="danger">Failed to login because {this.state.err}</Alert>)
            }
        }

        let signupLink = null;
        if (!this.props.isAdminPage) {
            signupLink=(
                <div>
                    <h4>Don't have an account?</h4>
                    <Link to={(this.props.isAdminPage) ? signupPath : signupPathWeb}>
                        <Button>Signup</Button>
                    </Link>
                </div>
            );
        }
        return (
            <Container>
                <Form onSubmit={this.onSubmit}>
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

const LoginPageComponent = connect(
    null,
    (dispatch) => {
        return {
            didAuthenticate: (state) => {dispatch(loginAction(state))}
        };
    }
)(CookieAwareComponent(withCookies(LoginPage)));

export default LoginPageComponent;

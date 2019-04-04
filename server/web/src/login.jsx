import React from 'react';
import {Redirect} from 'react-router-dom';
import { Button, Container, FormGroup, FormControl, ControlLabel, Alert, Form } from "react-bootstrap";
import {loginUrl} from './config.js';
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {landingPath} from './routes.js';

const LOGIN_ACTION = 'LOGIN';

const initialState = {
  isAuthenticated: false
}

export function loginAction() {
    return {type: LOGIN_ACTION};
}

export function loginReducer(state = initialState, action) {
    switch(action.type) {
        case LOGIN_ACTION:
            return Object.assign({}, state, {isAuthenticated: true});
        default:
            return state;
    }
}

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
        fetch(loginUrl, {
            method: 'POST',
            body: JSON.stringify({
                email: this.state.email,
                password: this.state.password
            })
        })
            .then(response => response.json())
            .then((data) => {
                console.log(data.Result.sessionId);
                cookies.set('sessionId', data.Result.sessionId);
                this.props.didAuthenticate();
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
        let { from } = this.props.location.state || { from: { pathname: landingPath } };

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
            </Container>
        );
    }
}

const LoginPageComponent = connect(
    null,
    (dispatch) => {
        return {
            didAuthenticate: () => {dispatch(loginAction())}
        };
    }
)(CookieAwareComponent(withCookies(LoginPage)));

export default LoginPageComponent;

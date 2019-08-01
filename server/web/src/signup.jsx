import React from 'react';
import { Button, Container, Alert, Form } from 'react-bootstrap';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {landingPathWeb} from './routes.js';
import RedirectPage from './redirect-flash.jsx';
import apiServiceConnect from './api/api_service_connect';

const SIGNUP_ACTION = 'SIGNUP';
const SIGNUP_REDIRECT_DURATION_SECONDS = 10;
const SIGNUP_REDIRECT_MESSAGE = "Thank you for signing up. Please login using your newly created account.";

// Unspecified gender is default 3 on server side.
// TODO: need a better way to keep these in sync.
const DEFAULT_GENDER = 3;

const initialState = {
  didSignup: false
}

export function signupAction() {
    return {type: SIGNUP_ACTION};
}

export function signupReducer(state = initialState, action) {
    switch(action.type) {
        case SIGNUP_ACTION:
            return Object.assign({}, state, {didSignup: true});
        default:
            return state;
    }
}

export class SignupPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            redirectToReferrer: false,
            firstName: "",
            lastName: "",
            email: "",
            gender: DEFAULT_GENDER,
            birthdate: "",
            phoneNumber: "",
            password: "",
            confirmPassword: ""
        };

        this.handleChange = this.handleChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.handleChangeInt = this.handleChangeInt.bind(this);
    }

    validateForm() {
        return this.state.firstName.length > 0 
            && this.state.lastName.length > 0 
            && this.state.email.length > 0 
            && this.state.birthdate.length > 0 
            && this.state.phoneNumber.length > 0 
            && this.state.password.length > 0
            && this.state.password == this.state.confirmPassword 
            ;
    }

    handleChange(event) {
        this.setState({
            [event.target.id]: event.target.value
        });
    }

    handleChangeInt(event) {
        this.setState({
            [event.target.id]: parseInt(event.target.value, 10)
        });
    }

    onSubmit(event) {
        const {cookies} = this.props;
        event.preventDefault();
        // send to api server
        this.props.apiService.signup(
            this.state.firstName, 
            this.state.lastName,
            this.state.email,
            this.state.gender,
            this.state.birthdate,
            this.state.phoneNumber,
            this.state.password
        )
            .then((data) => {
                this.props.didSignup();
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
        let { from } = this.props.location.state || { from: { pathname: landingPathWeb } };

        if (!!redirectToReferrer) {
            return <RedirectPage to={from} duration={SIGNUP_REDIRECT_DURATION_SECONDS} message={SIGNUP_REDIRECT_MESSAGE}/>;
        }

        let alert;
        if (this.state.submitState) {
            if (this.state.submitState === "SUCCESS") {
                alert = (<Alert variant="success">Succesfully signed up.</Alert>);
            } else if (this.state.submitState === "ERROR") {
                alert = (<Alert variant="danger">Failed to signup because {this.state.err}</Alert>)
            }
        }
        return (
            <Container className="panel-body">
                <Form onSubmit={this.onSubmit}>
                    <Form.Group controlId="firstName">
                        <Form.Label>First Name</Form.Label>
                        <Form.Control
                            autoFocus
                            type="text"
                            value={this.state.firstName}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="lastName">
                        <Form.Label>Last Name</Form.Label>
                        <Form.Control
                            autoFocus
                            type="text"
                            value={this.state.lastName}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="email">
                        <Form.Label>Email</Form.Label>
                        <Form.Control
                            autoFocus
                            type="email"
                            value={this.state.email}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="phoneNumber">
                        <Form.Label>Phone Number</Form.Label>
                        <Form.Control
                            autoFocus
                            type="text"
                            value={this.state.phoneNumber}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="birthdate">
                        <Form.Label>Birthdate (yyyy-mm-dd)</Form.Label>
                        <Form.Control
                            autoFocus
                            type="text"
                            value={this.state.birthdate}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="gender">
                        <Form.Label>Gender</Form.Label>
                        <Form.Control 
                            onChange={this.handleChangeInt}
                            as="select">
                            <option value={3}>Gender</option>
                            <option value={2}>Male</option>
                            <option value={1}>Female</option>
                        </Form.Control>
                    </Form.Group>
                    <Form.Group controlId="password">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            value={this.state.password}
                            onChange={this.handleChange}
                            type="password"
                        />
                    </Form.Group>
                    <Form.Group controlId="confirmPassword">
                        <Form.Label>Confirm Password</Form.Label>
                        <Form.Control
                            value={this.state.confirmpassword}
                            onChange={this.handleChange}
                            type="password"
                        />
                    </Form.Group>
                    <Button
                        block
                        disabled={!this.validateForm()}
                        type="submit"
                    >
                        Signup
                    </Button>
                </Form>
                <div className="message-container">
                    {alert}
                </div>
            </Container>
        );
    }
}

const SignupPageComponent = apiServiceConnect(
    null,
    (dispatch) => {
        return {
            didSignup: () => {dispatch(signupAction())}
        };
    }
)(CookieAwareComponent(withCookies(SignupPage)));

export default SignupPageComponent;

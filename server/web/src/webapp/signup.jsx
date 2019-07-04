import React from 'react';
import {Redirect} from 'react-router-dom';
import { Button, Container, FormGroup, FormControl, ControlLabel, Alert, Form } from 'react-bootstrap';
import {signupUrl} from '../config.js';
import { connect } from 'react-redux';
import CookieAwareComponent from '../cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {landingPath} from '../routes.js';

const SIGNUP_ACTION = 'SIGNUP';

// unspecified
const DEFAULT_GENDER = 3;

const initialState = {
  isAuthenticated: false
}

export function signupAction() {
    return {type: LOGIN_ACTION};
}

export function signupReducer(state = initialState, action) {
    switch(action.type) {
        case SIGNUP_ACTION:
            return Object.assign({}, state, {isAuthenticated: true});
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
        fetch(signupUrl, {
            method: 'POST',
            body: JSON.stringify({
                firstName: this.state.firstName,
                lastName: this.state.lastName,
                email: this.state.email,
                gender: this.state.gender,
                birthdate: this.state.birthdate,
                phoneNumber: this.state.phoneNumber,
                password: this.state.password
            })
        })
            .then(response => response.json())
            .then((data) => {
                console.log(data);
                // cookies.set('sessionId', data.Result.sessionId);
                // this.props.didAuthenticate();
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
                alert = (<Alert variant="success">Succesfully signed up.</Alert>);
            } else if (this.state.submitState === "ERROR") {
                alert = (<Alert variant="danger">Failed to signup because {this.state.err}</Alert>)
            }
        }
        return (
            <Container>
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

const SignupPageComponent = connect(
    null,
    (dispatch) => {
        return {
            didAuthenticate: () => {dispatch(signupAction())}
        };
    }
)(CookieAwareComponent(withCookies(SignupPage)));

export default SignupPageComponent;

import React from 'react';
import '../scss/group_register_page.scss';
import {Alert, Button, Container, Col, Form, Row} from 'react-bootstrap';
import { Redirect } from 'react-router-dom';
import ClipLoader from 'react-spinners/ClipLoader';

import apiServiceConnect from '../api/api_service_connect';
import {sendVerificationEmailApiModule} from '../api/send_verification_email_api_module';
import {getFetchInfo} from '../api/api_module';
import {landingPathWeb} from '../routes';


class VerifyEmailPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            email: "",
            sentEmail: "",
            isVerified: false,
        };

        this.handleChange = this.handleChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.onVerifiedClick = this.onVerifiedClick.bind(this);
    }

    validateForm() {
        return this.state.email.length;
    }

    handleChange(event) {
        this.setState({
            [event.target.id]: event.target.value
        });
    }

    onSubmit(event) {
        event.preventDefault();
        const { email } = this.state;
        this.setState({sentEmail: email});
        this.props.sendVerificationEmail(email);
    }

    onVerifiedClick() {
        this.setState({isVerified: true})
    }

    renderForm() {
        return (
            <Form onSubmit={this.onSubmit}>
                <Form.Group controlId="email">
                    <Form.Label>Waterloo Email</Form.Label>
                    <Form.Control
                        autoFocus
                        type="text"
                        value={this.state.email}
                        onChange={this.handleChange}
                    />
                </Form.Group>
                <Button
                    block
                    disabled={!this.validateForm()}
                    type="submit"
                >
                    Send Verification Email
                </Button>
            </Form>
        );
    }

    renderPostEmailMessaging() {
        const { fetchState, errorMessage } = this.props.sendVerificationEmailFetchInfo;
        const email = this.state.sentEmail;
        let body = null;
        // TODO(wojtek): This flow is not the greatest, since if someone clicks the "I verified"
        // button, but they haven't verified, then it will just resend them back to the same
        // screen.
        if (fetchState === "success") {
            body = (
                <Alert variant="success">
                    Email verification link sent to "{email}". Click on the link to verify that you
                    are a University of Waterloo student. When you are done, come back to this page
                    and click "I Already Verified". If you didn't receive the email, try again, by
                    pressing the "Send Verification Email" button.
                </Alert>
            );
        } else if (fetchState === "error") {
            body = (
                <Alert variant="danger">
                    Failed to send email to "{email}" with error: {errorMessage}. Try again,
                    by pressing the "Send Verification Email" button.
                </Alert>
            );
        } else if (fetchState === "fetching") {
            body = <ClipLoader />
        }
        return body;
    }

    render() {
        if (this.state.isVerified) {
            return <Redirect to={landingPathWeb} />;
        }

        return (
            <Container>
                <Row>
                    <Col>
                        <h3>Verify Email</h3>
                        <p>
                          { "Hive is the mentorship platform for the University of Waterloo. As a result, " +
                            "you must have a uwaterloo.ca or edu.uwaterloo.ca email to take part in the program." }
                        </p>
                    </Col>
                </Row>
                <Row><Col>{ this.renderForm() }</Col></Row>
                <Row><Col>{ this.renderPostEmailMessaging() }</Col></Row>
                <Button
                    block
                    onClick={this.onVerifiedClick}
                >
                    I Already Verified
                </Button>
            </Container>
        );
    }
}

const VerifyEmailPageComponent = apiServiceConnect(
    (state) => {
        return {
            sendVerificationEmail: sendVerificationEmailApiModule.getData(state),
            sendVerificationEmailFetchInfo: getFetchInfo(sendVerificationEmailApiModule, state),
        }
    },
    (dispatch) => {
        return {
            sendVerificationEmail: (email) => {
              return dispatch(sendVerificationEmailApiModule.getApiExecuteAction({email}))
            },
        }
    }
)(VerifyEmailPage);

export default VerifyEmailPageComponent;

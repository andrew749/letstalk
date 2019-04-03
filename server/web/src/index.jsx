import ReactDOM from 'react-dom';
import React from 'react';
import { Button, FormGroup, FormControl, ControlLabel, Alert, Form } from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {HiveToolTemplate} from './hive_tool_template.jsx';
import {loginUrl} from './config.js';

export class LoginPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
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
                console.log(cookies)
                console.log(data.Result.sessionId);
                cookies.set('sessionId', data.Result.sessionId);
                this.setState({
                    submitState: 'SUCCESS'
                });
            }).catch(err => {
                this.setState({
                    submitState: 'ERROR',
                    err: err.body
                });
            });

    }

    render() {

        let alert;
        if (this.state.submitState) {
            if (this.state.submitState === "SUCCESS") {
                alert = (<Alert variant="success">Succesfully logged in.</Alert>);
            } else if (this.state.submitState === "ERROR") {
                alert = (<Alert variant="danger">Failed to login because {this.state.err}</Alert>)
            }
        }
        return (
            <HiveToolTemplate>
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
            </HiveToolTemplate>
        );
    }
}

const CookieAwareLoginPage = CookieAwareComponent(withCookies(LoginPage));

ReactDOM.render(<CookieAwareLoginPage/>, document.getElementById('content'));

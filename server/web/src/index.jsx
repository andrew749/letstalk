import ReactDOM from 'react-dom';
import React from 'react';
import { Button, FormGroup, FormControl, ControlLabel, Alert } from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx'
import {withCookies} from 'react-cookie';

const serverUrl = "http://localhost/v1";


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
        fetch(`${serverUrl}/login`, {
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
                alert = (<Alert bsStyle="success">Succesfully logged in.</Alert>);
            } else if (this.state.submitState === "ERROR") {
                alert = (<Alert bsStyle="danger">Failed to login because {this.state.err}</Alert>)
            }
        }
        return (
            <div>
                <h2>Hive</h2>
                <form onSubmit={this.onSubmit}>
                    <FormGroup controlId="email" bsSize="large">
                        <ControlLabel>Email</ControlLabel>
                        <FormControl
                            autoFocus
                            type="email"
                            value={this.state.email}
                            onChange={this.handleChange}
                        />
                    </FormGroup>
                    <FormGroup controlId="password" bsSize="large">
                        <ControlLabel>Password</ControlLabel>
                        <FormControl
                            value={this.state.password}
                            onChange={this.handleChange}
                            type="password"
                        />
                    </FormGroup>
                    <Button
                        block
                        bsSize="large"
                        disabled={!this.validateForm()}
                        type="submit"
                    >
                        Login
                    </Button>
                </form>
                <div className="message-container">
                    {alert}
                </div>
            </div>
        );
    }
}

const CookieAwareLoginPage = CookieAwareComponent(withCookies(LoginPage));

ReactDOM.render(<CookieAwareLoginPage/>, document.getElementById('content'));

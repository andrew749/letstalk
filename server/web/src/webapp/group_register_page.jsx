import React from 'react';
import { Container } from "react-bootstrap";
import {Alert, Spinner} from 'react-bootstrap';

import apiServiceConnect from '../api/api_service_connect';

import { landingPathWeb } from '../routes';
import TimedRedirect from '../redirect-flash';

const REGISTER_ACTION = "REGISTERING";
const REGISTER_SUCCESSFULLY = "REGISTER_SUCCESS";
const REGISTER_ERROR = "REGISTER_ERROR";

const initialState = {
    registering: false,
    registeredSuccessfully: false,
    registerError: undefined
};

function getRegisterAction() {
    return {type: REGISTER_ACTION};
}

function getRegisterSuccessfullyAction() {
    return {type: REGISTER_SUCCESSFULLY};
}

function getRegisterErrorAction(message) {
    return {type: REGISTER_ERROR, registerError: message};
}

export function groupRegisterReducer(state = initialState, action) {
    switch(action.type) {
        case REGISTER_ACTION:
            return Object.assign({}, state, {registering: true});
        case REGISTER_SUCCESSFULLY:
            return Object.assign({}, state, {registering: false, registeredSuccessfully: true});
        case REGISTER_ERROR:
            return Object.assign({}, state, {registering: false, registeredSuccessfully: false, registerError: action.registerError});
        default:
            return state;
    }
}

class GroupRegisterPage extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.props.apiService.enrollInGroup(
            this.uuid(), 
            this.props.beginRegisteringCallback, 
            this.props.registeredSuccessfullyCallback, 
            this.props.registerErrorCallback
        );
    }

    uuid() {
        return this.props.match.params[0];
    }

    render() {
        let body;
        if (!!this.props.registering) {
            body = (
                <Spinner width={30} height={30} animation="border"/>
            );
        } else if (!!this.props.registeredSuccessfully) {
            body = (
                <div>
                    <Alert variant="success">Successfully registered with group</Alert>
                    <TimedRedirect to={landingPathWeb} duration={3} message={"Redirecting you soon..."}></TimedRedirect>
                </div>
            );
        } else if (!!this.props.registerError) {
            body = (
                <div>
                    <Alert variant="danger">Failed to register with group {JSON.stringify(this.props.registerError.message)}</Alert>
                    <p> Please try the link again or email <a href={`mailto:uw.hive@gmail.com`}>uw.hive@gmail.com</a></p>
                </div>
            );
        }

        return (
            <Container className="panel-body">
                <h3>Registering with group {this.uuid()}</h3>
                <Spinner animation="border"/>
                {body}
            </Container>
        );
    }
}

const GroupRegisterPageComponent = apiServiceConnect(
    (state) => {
        return {
            registering: state.groupRegisterReducer.registering,
            registeredSuccessfully: state.groupRegisterReducer.registeredSuccessfully,
            registerError: state.groupRegisterReducer.registerError
        }
    },
    (dispatch) => {
        return {
            beginRegisteringCallback: () => dispatch(getRegisterAction()),
            registeredSuccessfullyCallback: () => dispatch(getRegisterSuccessfullyAction()),
            registerErrorCallback: (err) => dispatch(getRegisterErrorAction(err))
        }
    }
)(GroupRegisterPage);

export default GroupRegisterPageComponent;
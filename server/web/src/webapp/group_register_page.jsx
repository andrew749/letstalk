import React from 'react';
import '../scss/group_register_page.scss';
import {Alert, Container} from 'react-bootstrap';
import ClipLoader from 'react-spinners/ClipLoader';

import apiServiceConnect from '../api/api_service_connect';


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
            return Object.assign({}, state, {registering: false, registeredSuccessfully: false,
                                             registerError: action.registerError.serverMessage});
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
                <div>
                    <ClipLoader />
                </div>
            );
        } else if (!!this.props.registeredSuccessfully) {
            body = (
                <Alert variant="success">Successfully joined group with id "{this.uuid()}"</Alert>
            );
        } else if (!!this.props.registerError) {
            body = (
                <Alert variant="danger">
                    Failed to join group with id "{this.uuid()}" with error: {this.props.registerError}
                </Alert>
            );
        }

        return (
            <Container>
                <h3>Registering with group "{this.uuid()}"</h3>
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

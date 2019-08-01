import React from 'react';
import { Button, Modal } from 'react-bootstrap';
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';

const SHOW_ACTION = 'SHOW';
const HIDE_ACTION = 'HIDE';

const initialState = {
  isVisible: false
}

export function showAction() {
    return {type: SHOW_ACTION};
}

export function hideAction() {
    return  {type: HIDE_ACTION};
}

export function modalReducer(state = initialState, action) {
    switch(action.type) {
        case SHOW_ACTION:
            return Object.assign({}, state, {isVisible: true});
        case HIDE_ACTION:
            return Object.assign({}, state, {isVisible: false});
        default:
            return state;
    }
}

export class ModalContainer extends React.Component {
  render() {
    return (
      <Modal size="lg" centered show={this.props.isVisible} onHide={this.props.hideModal}>
        <Modal.Header closeButton>
          <Modal.Title>Modal heading</Modal.Title>
        </Modal.Header>
        <Modal.Body>Woohoo, you're reading this text in a modal!</Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={this.props.hideModal}>
            Close
          </Button>
          <Button variant="primary" onClick={this.props.hideModal}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>
    );
  }
}

const ModalContainerComponent = connect(
  (state) => state.modalReducer,
  (dispatch) => {
    return {
        hideModal: (state) => {dispatch(hideAction(state))}
    };
  }
)(CookieAwareComponent(withCookies(ModalContainer)));

export default ModalContainerComponent;
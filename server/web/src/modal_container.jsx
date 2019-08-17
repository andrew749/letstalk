import React from 'react';
import { Form, Button, Modal, Row, Col } from 'react-bootstrap';
import CookieAwareComponent from './cookie_aware_component.jsx';
import { withCookies } from 'react-cookie';
import { onChange } from './util.js';
import apiServiceConnect from './api/api_service_connect';
import {fetchGroupsApiModule} from './api/fetch_groups'
import {getCurrentGroup} from './group_context_reducer';


const SHOW_ACTION = 'SHOW';
const HIDE_ACTION = 'HIDE';
export const MODAL_TYPES = {
  ADD_MEMBER: "add_member",
  DELETE_MEMBER: "delete_member",
  CREATE_MATCHING_ROUND: "create_matching_round",
  COMMIT_MATCHING_ROUND: "commit_matching_round",
  NONE: "none"
};

const initialState = {
  isVisible: false,
  modalType: MODAL_TYPES.NONE
}

export function showAction(modalType) {
    return {type: SHOW_ACTION, modalType: modalType};
}

export function hideAction() {
    return  {type: HIDE_ACTION};
}

export function modalReducer(state = initialState, action) {
    switch(action.type) {
        case SHOW_ACTION:
            return Object.assign({}, state, {isVisible: true, modalType: action.modalType});
        case HIDE_ACTION:
            return Object.assign({}, state, {isVisible: false});
        // case SHOW_ADD_MEMBER_ACTION:
        //     return Object.assign({}, state, {isVisible: true, modalType: MODAL_TYPES.ADD_MEMBER});
        default:
            return state;
    }
}

export class ModalContainer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      round: 1,
      createMatchingRoundModel: {
          matchingRoundName: "",
          maxLowerYearsPerUpperYear: 0,
          maxUpperYearsPerLowerYear: 0,
          youngestUpperGradYear: 0,
      }
    };
    this.onChange = onChange.bind(this, 'createMatchingRoundModel');
    this.saveMatchingRound = this.saveMatchingRound.bind(this);
  }

  saveMatchingRound(e) {
    e.preventDefault();
    e.stopPropagation();
    const model = this.state.createMatchingRoundModel;
    this.props.hideModal();
    // this.props.apiService.createManagedGroup(model.matchingRoundName, () => {/* When we start creating */}, () => {this.props.refreshPage()}, (err) => {console.warn(err);});
  }

  render() {
    let modalHeaderText = "Modal Heading";
    let modalBody = (<div>Modal Body Text</div>);
    let modalFooter = "Modal Footer";

    switch (this.props.modalType) {
      case MODAL_TYPES.NONE:
        break;
      case MODAL_TYPES.ADD_MEMBER:
        modalHeaderText = "Add Members"; 
        // TODO: When group selector is it's own component with reducer, we need to change the link
        // to whichever group is currently selected
        if (this.props.groups && this.props.groups.length == 0) {
          modalBody = (<div>{"Please create a group first!"}</div>);
        } else if (!this.props.selectedGroup) {
          modalBody = (<div>{"Please select a group first!"}</div>);
        } else {
          modalBody = (
            <div>
              {"Copy the following link into your email communications to allow your members to be added to your group: "}
              <a href={`mailto:${this.props.selectedGroup.managedGroupReferralEmail}`}>
                {this.props.selectedGroup.managedGroupReferralEmail}
              </a>
            </div>
          );
        }
        
        modalFooter = "";
        break;
      case MODAL_TYPES.DELETE_MEMBER:
          modalHeaderText = "Delete Members"; 
          // TODO: When group selector is it's own component with reducer, we need to change the link
          // to whichever group is currently selected
          modalBody = (
            <p>Are you sure you want to delete the selected users?</p>
          );
          modalFooter = "";
          break;
      case MODAL_TYPES.CREATE_MATCHING_ROUND:
          modalHeaderText = "Create New Matching Round";
          modalBody = (
            <div>
              <Form onSubmit={this.saveMatchingRound}>
                <Form.Group controlId="formCreateMatchingRound">
                  <Row className="matchround-form-row">
                    <Form.Label column sm="3">
                      Round Name
                    </Form.Label>
                    <Col sm="9">
                      <Form.Control size="lg" type="text" name="matchingRoundName" placeholder={"Matching Round " + this.state.round} onChange={this.onChange} />
                    </Col>
                  </Row>
                  <Row className="matchround-form-row">
                    <Form.Label column sm="3">
                      Max Mentees/Mentor
                    </Form.Label>
                    <Col sm="9">
                      <Form.Control size="lg" type="number" name="maxLowerYearsPerUpperYear" placeholder={0} onChange={this.onChange} />
                    </Col>
                  </Row>
                  <Row className="matchround-form-row">
                    <Form.Label column sm="3">
                      Max Mentors/Mentee
                    </Form.Label>
                    <Col sm="9">
                      <Form.Control size="lg" type="number" name="maxUpperYearsPerLowerYear" placeholder={0} onChange={this.onChange} />
                    </Col>
                  </Row>
                  <Row className="matchround-form-row">
                    <Form.Label column sm="3">
                      Min Mentor Grad Year
                    </Form.Label>
                    <Col sm="9">
                      <Form.Control size="lg" type="number" name="youngestUpperGradYear" placeholder={0} onChange={this.onChange} />
                    </Col>
                  </Row>
                </Form.Group>
                <div className="text-right">
                  <Button variant="primary" type="submit">
                      Create Matching Group
                  </Button>
                </div>
              </Form>
            </div>
          );
      default:
    }
    return (
      <Modal size="lg" centered show={this.props.isVisible} onHide={this.props.hideModal}>
        <Modal.Header closeButton>
          <Modal.Title>{modalHeaderText}</Modal.Title>
        </Modal.Header>
        <Modal.Body>{modalBody}</Modal.Body>
        {/* <Modal.Footer>
          <Button variant="secondary" onClick={this.props.hideModal}>
            Close
          </Button>
          <Button variant="primary" onClick={this.props.hideModal}>
            Save Changes
          </Button>
        </Modal.Footer> */}
      </Modal>
    );
  }
}

const ModalContainerComponent = apiServiceConnect(
  (state) => {
    return {
      selectedGroup: getCurrentGroup(state),
      groups: fetchGroupsApiModule.getData(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [], 
      ...state.modalReducer
    };
  },
  (dispatch) => {
    return {
        hideModal: (state) => {dispatch(hideAction(state))}
    };
  }
)(CookieAwareComponent(withCookies(ModalContainer)));

export default ModalContainerComponent;
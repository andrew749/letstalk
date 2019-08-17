import React from 'react';
import { Form, Button, Modal, Row, Col, Table } from 'react-bootstrap';
import BootstrapTable from 'react-bootstrap-table-next';
import apiServiceConnect from './api/api_service_connect';
import CookieAwareComponent from './cookie_aware_component.jsx';
import { withCookies } from 'react-cookie';
import { onChange } from './util.js';
import {fetchGroupsApiModule} from './api/fetch_groups';
import {fetchMembersApiModule} from './api/fetch_members';
import {getCurrentGroup} from './group_context_reducer';
import {matchRoundApiModule, commitMatchRoundApiModule} from './api/match_round_api_module';

const SHOW_ACTION = 'SHOW';
const HIDE_ACTION = 'HIDE';
export const MODAL_TYPES = {
  ADD_MEMBER: "add_member",
  DELETE_MEMBER: "delete_member",
  CREATE_MATCHING_ROUND: "create_matching_round",
  ADD_MEMBERS_TO_ROUND: "add_members_to_round",
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
      },
      userIds: [],
      matches: [],
    };

    this.selectRowProp = {
      mode: 'checkbox',
      clickToSelect: true,
      onSelect: this.onRowSelect,
      onSelectAll: this.onSelectAll,
    };

    this.onChangeMatchingRound = onChange.bind(this, 'createMatchingRoundModel');
  }

  onRowSelect = ({ id }, isSelected) => {
    if (isSelected) {
      this.setState({
        userIds: [ ...this.state.userIds, id ],
      });
    } else {
      this.setState({
        userIds: [ ...this.state.userIds.filter(idToCompare => idToCompare != id)]
      }); 
    }
    return ;
  }

  onSelectAll = (isSelected, rows) => {
    if (isSelected) {
      let ids = rows.map(row => row.id)
      this.setState({
          userIds: [...ids],
      });
    } else {
        this.setState({
            userIds: [],
        });
    }
  }

  saveMatchingParams = (e) => {
    e.preventDefault();
    e.stopPropagation();
    const form = e.currentTarget;
    if (form.checkValidity() === true) {
      this.props.showModal(MODAL_TYPES.ADD_MEMBERS_TO_ROUND);
    }
  }

  selectMembers = (e) => {
    e.preventDefault();
    e.stopPropagation();
    this.props.createNewMatchingRoundForGroup(
      this.props.selectedGroup.groupId,
      this.state.userIds,
      parseInt(this.state.createMatchingRoundModel.maxLowerYearsPerUpperYear),
      parseInt(this.state.createMatchingRoundModel.maxUpperYearsPerLowerYear), 
      parseInt(this.state.createMatchingRoundModel.youngestUpperGradYear)
    );
    this.props.showModal(MODAL_TYPES.COMMIT_MATCHING_ROUND);
    return;
  }

  commitMatchingRound = (e) => {
    e.preventDefault();
    e.stopPropagation();
    this.props.commitMatchingRoundForGroup(this.props.matchRoundId);
    this.props.hideModal();
    return;
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
          if (this.props.groups && this.props.groups.length == 0) {
            modalBody = (<div>{"Please create a group first!"}</div>);
          } else if (!this.props.selectedGroup) {
            modalBody = (<div>{"Please select a group first!"}</div>);
          } else {
            modalBody = (
              <div>
                <Form onSubmit={this.saveMatchingParams}>
                  <Form.Group controlId="formCreateMatchingRound">
                    <Row className="matchround-form-row">
                      <Form.Label column sm="3">
                        Round Name
                      </Form.Label>
                      <Col sm="9">
                        <Form.Control size="lg" type="text" name="matchingRoundName" placeholder={"Matching Round " + this.state.round} onChange={this.onChangeMatchingRound} required/>
                      </Col>
                    </Row>
                    <Row className="matchround-form-row">
                      <Form.Label column sm="3">
                        Max Mentees/Mentor
                      </Form.Label>
                      <Col sm="9">
                        <Form.Control size="lg" type="number" name="maxLowerYearsPerUpperYear" placeholder={0} onChange={this.onChangeMatchingRound} required/>
                      </Col>
                    </Row>
                    <Row className="matchround-form-row">
                      <Form.Label column sm="3">
                        Max Mentors/Mentee
                      </Form.Label>
                      <Col sm="9">
                        <Form.Control size="lg" type="number" name="maxUpperYearsPerLowerYear" placeholder={0} onChange={this.onChangeMatchingRound} required/>
                      </Col>
                    </Row>
                    <Row className="matchround-form-row">
                      <Form.Label column sm="3">
                        Min Mentor Grad Year
                      </Form.Label>
                      <Col sm="9">
                        <Form.Control size="lg" type="number" name="youngestUpperGradYear" placeholder={0} onChange={this.onChangeMatchingRound} required/>
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
          }
          break;
      case MODAL_TYPES.ADD_MEMBERS_TO_ROUND:
          const membersData = this.props.members ?  this.props.members.map(groupMember => {
            return ({
                id: groupMember.user.userId,
                name: groupMember.user.firstName + " " + groupMember.user.lastName,
                status: groupMember.status,
                email: groupMember.email,
                programName: groupMember.cohort ? (groupMember.cohort.programName + " " + groupMember.cohort.gradYear) : "No cohort"
            });
          }): [];
          const columns = [
            {
              dataField: 'name',
              text: 'Name'
            },
            {
              dataField: 'email',
              text: 'Email'
            },
            {
              dataField: 'status',
              text: 'Status'
            },
            {
              dataField: 'programName',
              text: 'Program'
            }
          ];
          modalHeaderText = "Adding Members to Matching Round";
          modalBody = (
            <div>
              {this.state.userIds.length === 0 && <p style={{color: "red"}}>You must select some members!</p>}
              <div className="main-table-container">
                <BootstrapTable keyField='id' data={membersData} columns={columns} selectRow={this.selectRowProp}/>
              </div>
              <div className="text-right">
                <Button variant="primary" onClick={this.selectMembers}>
                    Next
                </Button>
              </div>
            </div>
          );
          break;
      case MODAL_TYPES.COMMIT_MATCHING_ROUND:
          modalHeaderText = "Confirm Matching Round";
          modalBody = (
            <div>
              <div>
                <h2>
                  Match round ID: {this.props.matchRoundId}
                </h2>
                <h2>
                  Number of matches created: {this.props.matches.length}
                </h2>
                <h3>
                  Matches: 
                </h3>
                <div className="main-table-container">
                  <Table striped bordered hover>
                    <thead>
                        <tr>
                            <th>Mentor Name</th>
                            <th>Mentor Email</th>
                            <th>Mentee Name</th>
                            <th>Mentee Email</th>
                        </tr>
                    </thead>
                    <tbody>
                      {this.props.matches.map(match =>  {
                        return (
                          <tr>
                            <td>{match.mentor.user.firstName + " " + match.mentor.user.lastName}</td>
                            <td>{match.mentor.email}</td>
                            <td>{match.mentee.user.firstName + " " + match.mentee.user.lastName}</td>
                            <td>{match.mentee.email}</td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </Table>
                </div>
              </div>
              <div className="text-right">
                <Button variant="primary" onClick={this.commitMatchingRound}>
                    Commit Round
                </Button>
              </div>
            </div>
          );
          break;
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
      groups: fetchGroupsApiModule.isFinished(state) ? fetchGroupsApiModule.getData(state).managedGroups: undefined || [], 
      matchRoundId: matchRoundApiModule.isFinished(state) ? matchRoundApiModule.getData(state).matchRoundId: undefined || [],
      matches: matchRoundApiModule.isFinished(state) ? matchRoundApiModule.getData(state).matches: undefined || [],
      members: fetchMembersApiModule.isFinished(state) ? fetchMembersApiModule.getData(state): undefined || [],
      ...state.modalReducer
    };
  },
  (dispatch) => {
    return {
        createNewMatchingRoundForGroup: (groupId, userIds, maxLowerYearsPerUpperYear, maxUpperYearsPerLowerYear, youngestUpperGradYear) => dispatch(matchRoundApiModule.getApiExecuteAction({groupId, userIds, maxLowerYearsPerUpperYear, maxUpperYearsPerLowerYear, youngestUpperGradYear})),
        commitMatchingRoundForGroup: (matchRoundId) => dispatch(commitMatchRoundApiModule.getApiExecuteAction({matchRoundId})),
        showModal: (modalType) => {dispatch(showAction(modalType))},
        hideModal: (state) => {dispatch(hideAction(state))}
    };
  }
)(CookieAwareComponent(withCookies(ModalContainer)));

export default ModalContainerComponent;
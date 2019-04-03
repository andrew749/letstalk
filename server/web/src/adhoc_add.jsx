import React from 'react';
import CookieAwareComponent from './cookie_aware_component.jsx'
import { withCookies } from 'react-cookie';
import { Form, Button, Alert, Container, Row, Col } from 'react-bootstrap';
import './scss/notification_console.scss';
import { createMentorshipFromEmails } from './admin_api_controller.js';
import { createMentorshipFromEmails, deleteUser as deleteUserApi } from './admin_api_controller.js';
import {connect} from 'react-redux';

class AdminPanel extends React.Component {

  constructor(props) {
    super(props)
    this.onUpdateData = this.onUpdateData.bind(this);
    this.state = {
      // to be sent to server to echo back with information
      notificationState: {},
      adhocMatchingToolModel: {
        mentorEmail: undefined,
        menteeEmail: undefined,
        error: undefined
      },
      userDeletionToolModel: {
        userId: undefined,
        firstName: undefined,
        lastName: undefined,
        email: undefined,
        error: undefined
      },
    }

    this.adhocMatchingToolChange = this.onChange.bind(this, 'adhocMatchingToolModel');
    this.userDeletionToolChange = this.onChange.bind(this, 'userDeletionToolModel');
    this.createConnection = this.createConnection.bind(this);
    this.deleteUser = this.deleteUser.bind(this);
  }

  onUpdateData(newNotificationData) {
    this.setState({ notificationState: newNotificationData });
  }

  createConnection(e) {
    e.preventDefault();
    e.stopPropagation();
    console.log("Create Connection");
    console.log(this.state);
    // validation
    if (!this.state.adhocMatchingToolModel.menteeEmail || !this.state.adhocMatchingToolModel.mentorEmail) {
      console.log("Not enough fields filled out");
      this.setState({
        adhocMatchingToolModel: {
          ...this.state.adhocMatchingToolModel,
          error: "Missing required field"
        }
      })
      return;
    }
    console.log("Sending request");
    const mentorEmail = this.state.adhocMatchingToolModel.mentorEmail;
    const menteeEmail = this.state.adhocMatchingToolModel.menteeEmail;
    createMentorshipFromEmails(mentorEmail, menteeEmail)
      .then((data) => {
        // handle success response
        this.setState({ adhocMatchingToolModel: { error: undefined } });
        console.log("Successfully created mentorship.");
      }).catch(err => {
        this.setState({ adhocMatchingToolModel: { error: err.message } });
        console.warn("Failed to create mentorship");
      });
  }

  deleteUser(e) {
    e.preventDefault();
    e.stopPropagation();
    const model = this.state.userDeletionToolModel;
    if (!(model.email && model.firstName && model.lastName && model.userId )) {
      console.log("Not enough fields filled out");
      this.setState({
        userDeletionToolModel: {
          ...this.state.userDeletionToolModel,
          error: "Missing required field"
        }
      })
      return;
    }
    deleteUserApi(model.userId, model.firstName, model.lastName, model.email)
      .then((data) => {
        // handle success response
        console.log(data);
        this.setState({ userDeletionToolModel: { error: undefined } });
        console.log("Successfully deleted user");
      }).catch(err => {
        this.setState({ userDeletionToolModel: { error: err.message } });
        console.warn("Failed to delete user");
      });
  }

  onChange(model, event) {
    console.log(model);
    let fieldName = event.target.name;
    let fieldValue = event.target.value;
    console.log(fieldName);
    console.log(fieldValue);
    this.setState(
      prevState => ({
        [model]: {
          ...prevState[model],
          [fieldName]: fieldValue
        }
      })
    );
  }

  render() {
    const { cookies } = this.props;
    const adhocMatchingToolError = (this.state.adhocMatchingToolModel.error)
      ? <Alert key="adhocMatchingToolResponse" variant="danger">{this.state.adhocMatchingToolModel.error}</Alert>
      : undefined;
    const deleteToolError = (this.state.userDeletionToolModel.error)
      ? <Alert key="deleteToolResponse" variant="danger">{this.state.userDeletionToolModel.error}</Alert>
      : undefined;
    return (
        <Container>
          <Row>
            <h1>
              {"Adhoc Matching Tool"}
            </h1>
          </Row>
          <Row>
            <p>
              Enter the mentor and mentee emails for the relationship to create.
            </p>
          </Row>
          <Row>
            <Col style={{ paddingLeft: 0, paddingRight: 0 }} lg="12">
              <Form onSubmit={this.createConnection}>
                <Form.Group controlId="formMentorEmail">
                  <Form.Label>Mentor Email</Form.Label>
                  <Form.Control size="lg" type="email" name="mentorEmail" placeholder="Mentor Email" onChange={this.adhocMatchingToolChange} />
                </Form.Group>
                <Form.Group controlId="formMenteeEmail">
                  <Form.Label>Mentee Email</Form.Label>
                  <Form.Control size="lg" type="email" name="menteeEmail" placeholder="Mentee Email" onChange={this.adhocMatchingToolChange} />
                </Form.Group>
                <Button variant="primary" type="submit">
                  Match
              </Button>
                {adhocMatchingToolError}
              </Form>
            </Col>
          </Row>
          <hr/>
          <Row>
            <h1>
              {"User Deletion tool"}
            </h1>
          </Row>
          <Row>
            <p>
              Enter the details of the user to delete.
            </p>
          </Row>
          <Row>
            <Col style={{ paddingLeft: 0, paddingRight: 0 }} lg="12">
              <Form onSubmit={this.deleteUser}>
                <Form.Group controlId="formDeleteUser">
                  <Form.Control size="lg" type="number" name="userId" placeholder="User Id" onChange={this.userDeletionToolChange} />
                  <Form.Control size="lg" type="email" name="email" placeholder="Email" onChange={this.userDeletionToolChange} />
                  <Form.Control size="lg" type="text" name="firstName" placeholder="First Name" onChange={this.userDeletionToolChange} />
                  <Form.Control size="lg" type="text" name="lastName" placeholder="Last Name" onChange={this.userDeletionToolChange} />
                </Form.Group>
                <Button variant="primary" type="submit">
                  Delete
              </Button>
                {deleteToolError}
              </Form>
            </Col>
          </Row>
          <hr/>
        </Container>
    );
  }
}

const AdminPanelComponent = connect(null, null)(CookieAwareComponent(withCookies(AdminPanel)));

export default AdminPanelComponent;

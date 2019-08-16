import React from 'react';
import CookieAwareComponent from './cookie_aware_component.jsx'
import { withCookies } from 'react-cookie';
import { Form, Button, Alert, Container, Row, Col } from 'react-bootstrap';
import './scss/notification_console.scss';
import {onChange} from './util.js';
import apiServiceConnect from './api/api_service_connect';

class AdminPanel extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      // to be sent to server to echo back with information
      notificationState: {},
      adhocMatchingToolModel: {
        mentorEmail: undefined,
        menteeEmail: undefined,
        error: undefined,
        success: false
      }
    }

    this.adhocMatchingToolChange = onChange.bind(this, 'adhocMatchingToolModel');
    this.createConnection = this.createConnection.bind(this);
  }

  createConnection(e) {
    e.preventDefault();
    e.stopPropagation();
    // validation
    if (!this.state.adhocMatchingToolModel.menteeEmail || !this.state.adhocMatchingToolModel.mentorEmail) {
      console.log("Not enough fields filled out");
      this.setState({
        adhocMatchingToolModel: {
          ...this.state.adhocMatchingToolModel,
          error: "Missing required field",
          success: false
        }
      })
      return;
    }
    console.log("Sending request");
    const mentorEmail = this.state.adhocMatchingToolModel.mentorEmail;
    const menteeEmail = this.state.adhocMatchingToolModel.menteeEmail;
    this.props.apiService.createMentorshipFromEmails(mentorEmail, menteeEmail)
      .then((data) => {
        // handle success response
        this.setState({ adhocMatchingToolModel: { error: undefined, success: true } });
        console.log("Successfully created mentorship.");
      }).catch(err => {
        this.setState({ adhocMatchingToolModel: { error: err.message, success: false } });
        console.warn("Failed to create mentorship");
      });
  }

  render() {
    const { cookies } = this.props;
    const adhocMatchingToolError = (this.state.adhocMatchingToolModel.error)
      ? <Alert key="adhocMatchingToolResponse" variant="danger">{this.state.adhocMatchingToolModel.error}</Alert>
      : undefined;
    const adhocMatchingToolSuccess = (this.state.adhocMatchingToolModel.success)
      ? <Alert key="adhocMatchingToolResponse" variant="success">Successfully created mentorship connection.</Alert>
      : undefined;
    return (
        <Container className="panel-body">
          <div style={{ paddingLeft: 30, paddingTop: 30}}>
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
                  {adhocMatchingToolSuccess}
                </Form>
              </Col>
            </Row>
          </div>
        </Container>
    );
  }
}

const AdminPanelComponent = apiServiceConnect(null, null)(CookieAwareComponent(withCookies(AdminPanel)));

export default AdminPanelComponent;

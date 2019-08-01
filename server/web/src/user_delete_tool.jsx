import React from 'react';
import {Alert, Container, Form, Col, Row, Button} from 'react-bootstrap';
import {onChange} from './util.js';
import CookieAwareComponent from './cookie_aware_component.jsx'
import { withCookies } from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';

class UserDeleteTool extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            userDeletionToolModel: {
                userId: undefined,
                firstName: undefined,
                lastName: undefined,
                email: undefined,
                error: undefined,
                success: false
            }
        };

        this.userDeletionToolChange = onChange.bind(this, 'userDeletionToolModel');
        this.deleteUser = this.deleteUser.bind(this);
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
                    error: "Missing required field",
                    success: false
                }
            })
            return;
        }

        this.props.apiService.deleteUser(model.userId, model.firstName, model.lastName, model.email)
            .then((data) => {
                // handle success response
                console.log(data);
                this.setState({ userDeletionToolModel: { error: undefined, success: true } });
                console.log("Successfully deleted user");
            }).catch(err => {
                this.setState({ userDeletionToolModel: { error: err.message, success: false } });
                console.warn("Failed to delete user");
            });
    }

    render() {
        const deleteToolError = (this.state.userDeletionToolModel.error)
          ? <Alert key="deleteToolResponse" variant="danger">{this.state.userDeletionToolModel.error}</Alert>
          : undefined;
        const deleteToolSuccess = (this.state.userDeletionToolModel.success)
            ? <Alert key="deleteToolResponse" variant="success">Successfully deleted user.</Alert>
            : undefined;
        return (
            <Container className="panel-body">
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
                            {deleteToolSuccess}
                        </Form>
                    </Col>
                </Row>
            </Container>
        );
    }
}
const UserDeleteToolComponent = apiServiceConnect(null, null)(CookieAwareComponent(withCookies(UserDeleteTool)));
export default UserDeleteToolComponent;

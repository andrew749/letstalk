import React from 'react';
import { Redirect, Link} from 'react-router-dom';
import { Button, ButtonToolbar, Dropdown, DropdownButton, Container, Table} from "react-bootstrap";
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import apiServiceConnect from './api/api_service_connect';

const GROUPS = ['Hello Kitty', 'My Little Unicorn', 'Black Mamba'];

/**
 * Props:
 *  - isAdminPage: determine whether this is the admin page
 */
export class MatchingPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {
        const dropdownItems = GROUPS.map((groupName, i) => <Dropdown.Item key={i} eventKey={i.toString()}> {groupName} </Dropdown.Item>)
        return (
            <Container className="panel-body">
                <div className="group-info">
                    <h2>You are currently managing: </h2>
                    <ButtonToolbar>
                        <DropdownButton
                            title='Your Groups'
                            variant='Primary'
                            id='managed-groups-dropdown'
                        >
                            {dropdownItems}
                        </DropdownButton>
                    </ButtonToolbar>
                </div>
                <div className="panel-content">
                    <ButtonToolbar>
                        <Button variant="primary" size="lg">New Matching Round</Button>
                        {/* <Button variant="primary" size="lg">Delete members</Button>
                        <Button variant="primary" size="lg">Notify members</Button> */}
                    </ButtonToolbar>
                    <h2 className="mt-3">Matching Rounds</h2>
                    <div className="members-table-container">
                    <Table striped bordered hover>
                            <thead>
                                <tr>
                                    <th>#</th>
                                    <th>Name</th>
                                    <th>Stage</th>
                                    <th>Progress</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td>1</td>
                                    <td>Mark</td>
                                    <td>Otto</td>
                                    <td>@mdo</td>
                                </tr>
                                <tr>
                                    <td>2</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                                <tr>
                                    <td>3</td>
                                    <td>Jacob</td>
                                    <td>Thornton</td>
                                    <td>@fat</td>
                                </tr>
                            </tbody>
                        </Table>
                    </div>
                </div>
            </Container>
        );
    }
}

const MatchingPageComponent = apiServiceConnect()(CookieAwareComponent(withCookies(MatchingPage)));

export default MatchingPageComponent;

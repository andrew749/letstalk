import React from 'react';
import {LinkContainer} from 'react-router-bootstrap'
import {Navbar, Nav} from 'react-bootstrap';
import '../../scss/hive_tool_template.scss';
import {landingPathWeb} from '../../routes.js';

const CustomNavItem = (path, title) =>
    <LinkContainer to={path}>
        <Nav.Link>{title}</Nav.Link>
    </LinkContainer>


export default class HiveHeader extends React.Component {
    render() {
        return (
            <Navbar bg="dark" variant="dark">
                <Navbar.Brand>
                    <img
                        id="hive-logo"
                        alt="Hive"
                        src={require('../../img/logo.png')}
                        width="30"
                        height="30"
                    />
                </Navbar.Brand>
                <Nav>
                    {CustomNavItem(landingPathWeb, "Home")}
                </Nav>
            </Navbar>
        );
    }
}

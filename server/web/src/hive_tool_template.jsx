import ReactDOM from 'react-dom';
import React from 'react';
import {LinkContainer} from 'react-router-bootstrap'
import {Navbar, Nav} from 'react-bootstrap';
import './scss/hive_tool_template.scss';
import {adhocAddToolPath, deleteUserToolPath} from './routes.js';

const CustomNavItem = (path, title) =>
    <LinkContainer to={path}>
        <Nav.Link>{title}</Nav.Link>
    </LinkContainer>


export default class HiveToolTemplate extends React.Component {
    render() {
        return (
            <Navbar bg="dark" variant="dark">
                <Navbar.Brand href="#home">
                    <img
                        id="hive-logo"
                        alt="Hive"
                        src={require('./img/logo.png')}
                        width="30"
                        height="30"
                    />
                    {' Hive Admin Console'}
                </Navbar.Brand>
                <Nav>
                    {CustomNavItem(adhocAddToolPath, "Adhoc Matching Tool")}
                    {CustomNavItem(deleteUserToolPath, "Delete User")}
                </Nav>
            </Navbar>
        );
    }
}

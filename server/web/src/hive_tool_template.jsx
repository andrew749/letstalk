import ReactDOM from 'react-dom';
import React from 'react';
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {LinkContainer} from 'react-router-bootstrap'
import {Container, Navbar, Nav} from 'react-bootstrap';
import './scss/hive_tool_template.scss';
import {adhocAddToolPath, deleteUserToolPath, membersPath, matchingPath} from './routes.js';
import {HiveApiService} from './api_controller.js';

const CustomNavItem = (path, title) =>
    <LinkContainer to={path}>
        <Nav.Link>{title}</Nav.Link>
    </LinkContainer>

function isAuthenticated(state) {
    return state.loginReducer.isAuthenticated;
}

class HiveToolTemplate extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    populateMe() {
        if (!!this.props.isAuthenticated) {
            HiveApiService.me(
                data => this.setState({me: data}),
                err => console.log
            );
        }
    }

    componentDidMount() {
        this.populateMe();
    }

    componentDidUpdate() {
        this.populateMe();
    }

    render() {
        return (
            <Container className="navbar">
                <Navbar className="flex-column">
                    <div className="account-info">
                        <Navbar.Brand bg="blue">
                            <img
                                id="hive-logo"
                                alt="Hive"
                                src={require('./img/logo.png')}
                                width="70"
                                height="70"
                            />
                        </Navbar.Brand>
                        <div className="flex-column">
                            {'Hive Admin Console'}
                            <div>
                            {!!this.props.isAuthenticated && this.state && this.state.me ?  
                                !!this.state.me.Result ? 
                                    this.state.me.Result.email : "Unknown Email" 
                                : "Not authenticated"}
                            </div>
                        </div>
                    </div>
                    <Navbar.Toggle aria-controls="hive-navbar-nav" />
                    <Navbar.Collapse id="hive-navbar-nav">
                        <Nav defaultActiveKey="/" className="flex-column">
                            {CustomNavItem(membersPath, "Members")}
                            {CustomNavItem(matchingPath, "Matching")}
                            {CustomNavItem(adhocAddToolPath, "Adhoc Matching Tool")}
                            {CustomNavItem(deleteUserToolPath, "Delete User")}
                        </Nav>
                    </Navbar.Collapse>
                </Navbar>
            </Container>
        );
    }
}

const HiveToolComponent = connect(
    (state) => {
        return {
            isAuthenticated: isAuthenticated(state)
        };
    },
    null
)(CookieAwareComponent(withCookies(HiveToolTemplate)));

export default HiveToolComponent;
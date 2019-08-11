import { Redirect } from 'react-router-dom';
import React from 'react';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {LinkContainer} from 'react-router-bootstrap'
import {Container, Navbar, Nav} from 'react-bootstrap';
import './scss/hive_tool_template.scss';
import {adhocAddToolPath, deleteUserToolPath, loginPath, logoutPath, membersPath, matchingPath, groupManagementToolPath} from './routes.js';
import {logoutAction} from './login';
import {shouldFetchProfileAction, isAuthenticated, getProfile} from './api/api_controller'
import apiServiceConnect from './api/api_service_connect';

class NavbarContainer extends React.Component {

    constructor(props) {
        super(props);

        this.state = {};
    }


    componentDidMount() {
        this.props.fetchProfile();
    }

    onLogout = event => {
        // console.log(this.props);
        const {cookies} = this.props;
        event.preventDefault();
        // send to api server

        this.props.apiService.logout()
            .then((data) => {
                cookies.set('sessionId', null);
                this.props.didLogout();
                this.setState({
                    submitState: 'SUCCESS',
                    redirectToLogin: true
                });
            }).catch(err => {
                this.setState({
                    submitState: 'ERROR',
                    err: err.body
                });
            });
    };

    render() {
        let { redirectToLogin } = this.state;
        
        if (!!redirectToLogin && this.props.isAuthenticated) {
            return <Redirect to={loginPath} />;
        }

        let alert;
        if (this.state.submitState) {
            if (this.state.submitState === "ERROR") {
                alert = (<Alert variant="danger">Failed to logout because {this.state.err}</Alert>)
            }
        }

        const CustomNavItem = (path, title) => {
            if (path === logoutPath) {
                return <a className="nav-link" onClick={this.onLogout}>
                    {title}
                </a>
            }
            return <LinkContainer to={path}>
                <Nav.Link>{title}</Nav.Link>
            </LinkContainer>
        }

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
                            {!!this.props.isAuthenticated && this.props.profile ?
                                !!this.props.profile ? 
                                    this.props.profile.email : "Unknown Email" 
                                : "Not authenticated"}
                            </div>
                            <Nav>
                                {!this.props.isAuthenticated && CustomNavItem(loginPath, "Login")}
                                {!!this.props.isAuthenticated && CustomNavItem(logoutPath, "Logout")}
                            </Nav>
                        </div>
                        <div className="message-container">
                            {alert}
                        </div>
                    </div>
                    <Navbar.Toggle aria-controls="hive-navbar-nav" />
                    <Navbar.Collapse id="hive-navbar-nav">
                        <Nav defaultActiveKey="/" className="flex-column">
                            {CustomNavItem(membersPath, "Members")}
                            {CustomNavItem(matchingPath, "Matching")}
                            {CustomNavItem(adhocAddToolPath, "Adhoc Matching Tool")}
                            {CustomNavItem(deleteUserToolPath, "Delete User")}
                            {CustomNavItem(groupManagementToolPath, "Manage Groups")}
                        </Nav>
                    </Navbar.Collapse>
                </Navbar>
            </Container>
        );
    }
}

const NavbarContainerComponent = apiServiceConnect(
    (state) => {
        return {
            profile: getProfile(state),
            isAuthenticated: isAuthenticated(state)
        };
    },
    (dispatch) => {
        return {
            didLogout: (state) => {dispatch(logoutAction(state))},
            fetchProfile: () => {dispatch(shouldFetchProfileAction())}
        };
    }
)(CookieAwareComponent(withCookies(NavbarContainer)));

export default NavbarContainerComponent;
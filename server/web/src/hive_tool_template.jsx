import ReactDOM from 'react-dom';
import { Redirect } from 'react-router-dom';
import React from 'react';
import { connect } from 'react-redux';
import CookieAwareComponent from './cookie_aware_component.jsx';
import {withCookies} from 'react-cookie';
import {LinkContainer} from 'react-router-bootstrap'
import {Container, Navbar, Nav} from 'react-bootstrap';
import './scss/hive_tool_template.scss';
import {adhocAddToolPath, deleteUserToolPath, loginPath, logoutPath, membersPath, matchingPath} from './routes.js';
import {logoutAction} from './login';
import {HiveApiService} from './api_controller.js';

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
                ({ Result }) => {
                    if (!this.state.me || this.state.me.userId != Result.userId ) {
                        this.setState({ me: {
                            userId: Result.userId,
                            firstName: Result.firstName,
                            lastName: Result.lastName,
                            email: Result.email
                        }});
                    }
                },
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

    onLogout = event => {
        // console.log(this.props);
        const {cookies} = this.props;
        event.preventDefault();
        // send to api server

        HiveApiService.logout()
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
                            {!!this.props.isAuthenticated && this.state.me ?
                                !!this.state.me ? 
                                    this.state.me.email : "Unknown Email" 
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
    (dispatch) => {
        return {
            didLogout: (state) => {dispatch(logoutAction(state))}
        };
    }
)(CookieAwareComponent(withCookies(HiveToolTemplate)));

export default HiveToolComponent;
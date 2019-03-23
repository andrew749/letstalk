import ReactDOM from 'react-dom';
import React from 'react';
import {Navbar} from 'react-bootstrap';
import './scss/hive_tool_template.scss';

export class HiveToolTemplate extends React.Component {
    render() {
        return (
            <div>
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
                </Navbar>
                <div className="content-container">
                    { this.props.children }
                </div>
            </div>
        );
    }
}

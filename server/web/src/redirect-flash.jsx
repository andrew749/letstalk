import React from 'react';
import {Redirect} from 'react-router';

/**
 * Expects the following props:
 *  to: the page to redirect to after time
 *  message: a message to flash to the user
 *  duration: how long to keep the message before changing pages 
 */
export default class RedirectPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            redirect: false
        };
    }

    componentDidMount() {
        setTimeout(() => {
            this.setState({
                redirect: true,
            })
        }, this.props.duration * 1000);
    }

    render() {
        if (!!this.state.redirect) {
            return <Redirect to={this.props.to}/>;
        }
        return (
            <div>
                <p>{this.props.message}</p>
            </div>
        );
    }
}
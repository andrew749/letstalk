import React from 'react';
import {Redirect} from 'react-router-dom';

import { connect } from 'react-redux';

/**
 * Creates a redirect to the correct route based on the state of the redirectPathReducer.
 *
 *  redirectPath: the path to redirect a user to after authentication, if they entered the website
 *    from some other page that required auth.
 *  defaultPath: the default redirect path for this app (different between
 */
class PostAuthRedirect extends React.Component {
    redirectPath() {
        const { redirectPath, defaultPath } = this.props;
        if (!!redirectPath) {
            return redirectPath;
        } else {
            return defaultPath;
        }
    }

    render() {
        const pathname = this.redirectPath();
        const redirect = { pathname };
        return <Redirect to={redirect} />;
    }
}

export default connect(
    (state) => {
        return {
            redirectPath: state.postAuthReducer.redirectPath
        };
    },
    null,
)(PostAuthRedirect);

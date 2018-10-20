import React from 'react';
import {CookiesProvider} from 'react-cookie';

export default function CookieAwareComponent(WrappedComponent) {
    return class extends React.Component {
        render() {
            return (
                <CookiesProvider>
                    <WrappedComponent {...this.props}/>
                </CookiesProvider>
            );
        }
    }
}

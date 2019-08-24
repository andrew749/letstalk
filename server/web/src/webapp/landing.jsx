import React from 'react';
import { withCookies } from 'react-cookie';
import { Button, Container } from "react-bootstrap";
import apiServiceConnect from '../api/api_service_connect';

import CookieAwareComponent from '../cookie_aware_component.jsx';
import HiveLogoWithTitle from '../img/logo_with_title.png';


export class LandingPage extends React.Component {
    render() {
        return (
            <Container className="panel-body">
                <a href={"https://www.hiveapp.org"}>
                    <img id="hive-logo" src={HiveLogoWithTitle}/>
                </a>
                <p className="landing-text"> Thanks for signing up! We're currently building out our web experience. In the meantime, please go to <a href={"https://www.hiveapp.org"}>
                    our landing page
                </a>  or click the links below to install our app to fill out your profile and continue onboarding!</p>
                <div>
                    <a href={"https://itunes.apple.com/us/app/uw-hive/id1400424718?mt=8"}><img className="apple-link app-store-link" src={'https://linkmaker.itunes.apple.com/en-us/badge-lrg.svg?releaseDate=2018-09-28&kind=iossoftware&bubble=ios_apps'}/></a>
                    <a href={'https://play.google.com/store/apps/details?id=org.hiveapp.hive&pcampaignid=MKT-Other-global-all-co-prtnr-py-PartBadge-Mar2515-1'} ><img className={"google-link"} src={'https://play.google.com/intl/en_us/badges/images/generic/en_badge_web_generic.png'}/></a>
                </div>
            </Container>
        );
    }
}

const LandingPageComponent = apiServiceConnect(
    null,
    (dispatch) => {
        return {};
    }
)(CookieAwareComponent(withCookies(LandingPage)));

export default LandingPageComponent;

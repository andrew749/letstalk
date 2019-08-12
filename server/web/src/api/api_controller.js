import { loginUrl, logoutUrl, meUrl, signupUrl, mentorshipUrl, deleteUrl, getGroupMembersUrlBase, getManagedGroupsUrl, createNewManagedGroupUrl, registerWithManagedGroupUrl, getMatchRoundsUrl, createMatchRoundsUrl, userGroupUrl } from '../config.js'
import axios from 'axios';
import Cookies from 'universal-cookie';

const CREATE_MENTORSHIP_TYPE_DRY_RUN = "DRY_RUN";
const CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN = "NOT_DRY_RUN";

export const DID_AUTHENTICATE_ACTION = "DID_AUTHENTICATE";
export const AUTH_EXPIRED = "AUTH_EXPIRED";
export const FETCH_PROFILE = "FETCH_PROFILE";
export const FETCHING_PROFILE = "FETCHING_PROFILE";
export const DID_FETCH_PROFILE = "DID_FETCH_PROFILE";
export const FETCH_PROFILE_ERROR = "FETCH_PROFILE_ERROR";

// User did authenticate
export function didAuthenticateAction(sessionId) {
    return {type: DID_AUTHENTICATE_ACTION, sessionId: sessionId}
}

export function shouldFetchProfileAction() {
    return {type: FETCH_PROFILE};
}

export function fetchingProfileAction() {
    return {type: FETCHING_PROFILE};
}

export function didFetchProfileAction(profile) {
    return {type: DID_FETCH_PROFILE, profile: profile};
}

export function fetchProfileErrorAction(error) {
    return {type: FETCH_PROFILE_ERROR, error: error};
}

// Invalidate authentication
export function authExpiredAction() {
    return {type: AUTH_EXPIRED};
}

export function getShouldFetchProfile(state) {
    return state.apiServiceReducer.shouldFetchProfile;
}

export function isAuthenticated(state) {
    return state.apiServiceReducer.isValid;
}

export function getProfile(state) {
    return state.apiServiceReducer.profile;
}

// base state
const initialState = {
    sessionId: undefined,
    isValid: false,
    shouldFetchProfile: false,
};

export function apiServiceReducer(state = initialState, action) {
    switch(action.type) {
        case FETCH_PROFILE:
            return Object.assign({}, state, {shouldFetchProfile: true});
        case FETCHING_PROFILE:
            return Object.assign({}, state, {shouldFetchProfile: false});
        case DID_FETCH_PROFILE:
            return Object.assign({}, state, {shouldFetchProfile: false, profile: action.profile});
        case FETCH_PROFILE_ERROR:
            return Object.assign({}, state, {shouldFetchProfile: false, fetchProfileError: action.error});
        case DID_AUTHENTICATE_ACTION:
            return Object.assign({}, state, {isValid: true, sessionId: action.sessionId})
        case AUTH_EXPIRED:
            return Object.assign({}, state, {isValid: false, sessionId: undefined})
        default:
            return state;
    }
}

function getLocalStateForComponent(state) {
    return state.apiServiceReducer;
}

/**
 * A singleton that is used to make api calls to the hive webapp.
 */
export const HiveApiService = ((state, dispatch) => {
    
    const apiService = () => HiveApiService(state, dispatch);
    // holds internal state for the service.
    return {
        isAuthenticated: () => {
            return getLocalStateForComponent(state).isValid;
        },

        setSessionId: (sessionId) => {
            (new Cookies()).set('sessionId', sessionId);
            dispatch(shouldFetchProfileAction());
            dispatch(didAuthenticateAction(sessionId));
        },

        getSessionId: () => {
            return getLocalStateForComponent(state).sessionId;
        },

        /**
         * headers: {undefined || Headers} headers to send the request with.
         */
        hiveFetch: (url, method, body) => {
            let headers = {"sessionId": apiService().getSessionId(), "Content-Type": "application/json"};
            return axios({
                url:  url,
                method: method,
                headers: headers,
                data: body
            }).then(response => response.data)
                .catch(err => {
                    if (err.response && err.response.status === 401) {
                        dispatch(authExpiredAction());
                    }
                    throw err;
                });
        },

        signup: (firstName, lastName, email, gender, birthdate, phoneNumber, password) => {
            return apiService().hiveFetch(signupUrl, 'POST', {
                firstName: firstName,
                lastName: lastName,
                email: email,
                gender: gender,
                birthdate: birthdate,
                phoneNumber: phoneNumber,
                password: password
            });
        },

        login: (email, password) => {
            return apiService().hiveFetch(loginUrl, 'POST', {
                email: email,
                password: password
            }).then(data => {
                apiService().setSessionId(data.Result.sessionId);
                return data;
            });
        },

        logout: () => {
            return apiService().hiveFetch(logoutUrl, 'POST', undefined);
        },

        me: (started, done, error) => {
            started();
            return apiService().hiveFetch(meUrl, 'GET', undefined)
                .then((data) => done(data))
                .catch((err) => error(err));
        },

        /** 
         * createMentorshipFromEmails
         * 
         * Creates a new mentorship for the two specific users
         * 
         * Return:
         *  Promise for request
        */
        createMentorshipFromEmails: (mentorEmail, menteeEmail) => {
            return apiService().hiveFetch(mentorshipUrl, 'POST', {
                mentorEmail: mentorEmail,
                menteeEmail: menteeEmail,
                requestType: CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN
            });
        },

        /**
         * deleteUser
         * Delete the user with the specified parameters.
         * @param {*} userId 
         * @param {*} firstName 
         * @param {*} lastName 
         * @param {*} email 
         */
        deleteUser: (userId, firstName, lastName, email) => {
            return apiService().hiveFetch(deleteUrl, 'POST', {
                userId: parseInt(userId),
                firstName: firstName,
                lastName: lastName,
                email: email
            });
        },

        /**
         * fetchGroups Get groups an administrator manages.
         * @param {*} started Started fetching callback
         * @param {*} done Done fetching callback
         * @param {*} error Error fetching data
         */
        fetchGroups: (started, done, error) => {
            started();
            console.log("Fetching groups");
            return apiService().hiveFetch(getManagedGroupsUrl, 'GET', undefined)
                .then((data) => done(data))
                .catch((err) => error(err));
        },

        /**
         * fetchMembers Get members in a given group an administrator manages.
         * @param {*} started Started fetching callback
         * @param {*} done Done fetching callback
         * @param {*} error Error fetching data
         */
        fetchMembers: (groupId, started, done, error) => {
            if (!groupId) {
                groupId = 0;
            }
            let membersUrl = getGroupMembersUrlBase + groupId.toString();
            started();
            console.log("Fetching members");
            return apiService().hiveFetch(membersUrl, 'GET', undefined)
                .then((data) => done(data))
                .catch((err) => error(err));
        },

        /**
         * createManagedGroup Creates a new managed group
         * @param {*} groupName Name of the groupt to create
         */
        createManagedGroup: (groupName, started, done, error) => {
            started();
            console.log("Creating group " + groupName);
            return apiService().hiveFetch(createNewManagedGroupUrl, 'POST', {
                groupName: groupName,
            })
                .then((data) => done(data))
                .catch((err) => error(err));
        },

        /**
         * enrollInGroup
         */
        enrollInGroup: (uuid, started, done, error) => {
            started();
            return apiService().hiveFetch(registerWithManagedGroupUrl, 'POST', {
                groupUUID: uuid,
            })
                .then(done)
                .catch(error);
        },

        getMatchingRounds: (groupId, started, done, error) => {
            started();
            return apiService().hiveFetch(getMatchRoundsUrl + groupId, 'GET')
                .then(done)
                .catch(error);
        },

        createNewMatchingRound: (maxLowerYearsPerUpperYear, maxUpperYearsPerLowerYear, youngestUpperYearGrad, groupId, userIds, started, done, error) => {
            started();
            return apiService().hiveFetch(createMatchRoundsUrl, 'POST', {
                parameters: {
                    maxLowerYearsPerUpperYear,
                    maxUpperYearsPerLowerYear,
                    youngestUpperYearGrad, 
                },
                groupId: groupId,
                userIds: userIds,
            })
                .then(done)
                .catch(error);
        },
        deleteMemberFromGroup: (groupId, userId, started, done, error) => {
            started();
            return apiService().hiveFetch(userGroupUrl, 'DELETE', {
                groupId: groupId,
                userId: userId,
            })
                .then(done)
                .catch(error);
        }
    }
}

);

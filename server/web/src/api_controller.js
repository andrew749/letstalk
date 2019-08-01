import { loginUrl, logoutUrl, meUrl, signupUrl, mentorshipUrl, deleteUrl, getManagedGroupsUrl, createNewManagedGroupUrl, registerWithManagedGroupUrl } from './config.js'
import axios from 'axios';
import { registerWithGroupPathWeb } from './routes.js';

const CREATE_MENTORSHIP_TYPE_DRY_RUN = "DRY_RUN";
const CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN = "NOT_DRY_RUN";

export function fetchMiddleware(fetchPromise) {
    return new Promise((req, rej) => fetchPromise
        .then(response => response.data)
        .then((data) => {
            if (data.Error) {
                throw new Error(data.Error.message)
            }
            req(data);
        }).catch(err => rej(err)));
}

/**
 * A singleton that is used to make api calls to the hive webapp.
 */
export const HiveApiService = (() => {
    // holds internal state for the service.
    let instance = new Object();

    return {
        setSessionId: (sessionId) => {
            instance.sessionId = sessionId;
        },
        /**
         * headers: {undefined || Headers} headers to send the request with.
         */
        hiveFetch: (url, method, body) => {
            let headers = {"sessionId": instance.sessionId, "Content-Type": "application/json"};
            return fetchMiddleware(axios({
                url:  url,
                method: method,
                headers: headers,
                data: body
        }))},

        signup: (firstName, lastName, email, gender, birthdate, phoneNumber, password) => {
            return HiveApiService.hiveFetch(signupUrl, 'POST', {
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
            return HiveApiService.hiveFetch(loginUrl, 'POST', {
                email: email,
                password: password
            });
        },

        logout: () => {
            return HiveApiService.hiveFetch(logoutUrl, 'POST', undefined);
        },

        me: (done, error) => {
            return HiveApiService.hiveFetch(meUrl, 'GET', undefined)
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
            return HiveApiService.hiveFetch(mentorshipUrl, 'POST', {
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
            return HiveApiService.hiveFetch(deleteUrl, 'POST', {
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
            return HiveApiService.hiveFetch(getManagedGroupsUrl, 'GET', undefined)
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
            return HiveApiService.hiveFetch(createNewManagedGroupUrl, 'POST', {
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
            return HiveApiService.hiveFetch(registerWithManagedGroupUrl, 'POST', {
                groupUUID: uuid,
            })
                .then(done)
                .catch(error);
        }
    }
})();


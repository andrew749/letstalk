import { mentorshipUrl, deleteUrl } from './config.js'

const CREATE_MENTORSHIP_TYPE_DRY_RUN = "DRY_RUN";
const CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN = "NOT_DRY_RUN";

export function fetchMiddleware(fetchPromise) {
    return new Promise((req, rej) => fetchPromise
        .then(response => response.json())
        .then((data) => {
            if (data.Error) {
                throw new Error(data.Error.message)
            }
            req(data);
        }).catch(err => rej(err)));
}

/** 
 * createMentorshipFromEmails
 * 
 * Creates a new mentorship for the two specific users
 * 
 * Return:
 *  Promise for request
*/
export function createMentorshipFromEmails(mentorEmail, menteeEmail) {
    return fetchMiddleware(fetch(mentorshipUrl, {
        method: 'POST',
        body: JSON.stringify({
            mentorEmail: mentorEmail,
            menteeEmail: menteeEmail,
            requestType: CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN
        })
    }));
}

/**
 * deleteUser
 * Delete the user with the specified parameters.
 * @param {*} userId 
 * @param {*} firstName 
 * @param {*} lastName 
 * @param {*} email 
 */
export function deleteUser(userId, firstName, lastName, email) {
    return fetchMiddleware(fetch(deleteUrl, {
        method: 'POST',
        body: JSON.stringify({
            userId: parseInt(userId),
            firstName: firstName,
            lastName: lastName,
            email: email
        })
    }));
}
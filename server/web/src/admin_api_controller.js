import { mentorshipUrl } from './config.js'

const CREATE_MENTORSHIP_TYPE_DRY_RUN = "DRY_RUN";
const CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN = "NOT_DRY_RUN";

/** 
 * createMentorshipFromEmails
 * 
 * Creates a new mentorship for the two specific users
 * 
 * Return:
 *  Promise for request
*/
export function createMentorshipFromEmails(mentorEmail, menteeEmail) {
    return fetch(mentorshipUrl, {
        method: 'POST',
        body: JSON.stringify({
            mentorEmail: mentorEmail,
            menteeEmail: menteeEmail,
            requestType: CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN
        })
    });
}
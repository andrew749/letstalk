export const serverUrl = 'https://api.hiveapp.org';
//export const serverUrl = 'http://localhost'; // For dev
export const hiveDeepLinkRoot = 'hive:/';
//export const hiveDeepLinkRoot = 'exp://192.168.0.179:19000/--'; // For dev

export const apiUrl = `${serverUrl}/v1`;
export const adminUrl = `${serverUrl}/admin`

// api v1 endpoints
export const loginUrl = `${apiUrl}/login`;
export const signupUrl = `${apiUrl}/signup`;


// admin endpoints
export const mentorshipUrl = `${adminUrl}/mentorship`;
export const deleteUrl = `${adminUrl}/nuke_user`;
export const getManagedGroupsUrl = `${adminUrl}/get_managed_groups`;
export const createNewManagedGroupUrl = `${adminUrl}/create_managed_group`;
export const registerWithManagedGroupUrl = `${apiUrl}/enroll_managed_group`;

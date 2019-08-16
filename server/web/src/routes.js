const ADMIN_URL_PREFIX = "/admin_panel";

export const loginPath =`${ADMIN_URL_PREFIX}/login`;
export const logoutPath =`${ADMIN_URL_PREFIX}/logout`;
export const signupPath =`${ADMIN_URL_PREFIX}/signup`;
export const getManagedGroupsPath =`${ADMIN_URL_PREFIX}/get_managed_groups`;
export const membersPath =`${ADMIN_URL_PREFIX}/members`;
export const matchingPath =`${ADMIN_URL_PREFIX}/matching`;
export const adhocAddToolPath =`${ADMIN_URL_PREFIX}/adhoc_add`;
export const deleteUserToolPath =`${ADMIN_URL_PREFIX}/delete_user`;
export const groupManagementToolPath =`${ADMIN_URL_PREFIX}/manage_groups`;
export const landingPath = `${ADMIN_URL_PREFIX}/`;

const WEBAPP_URL_PREFIX = "/web";
export const signupPathWeb = `${WEBAPP_URL_PREFIX}/signup`;
export const loginPathWeb = `${WEBAPP_URL_PREFIX}/login`;
export const landingPathWeb = `${WEBAPP_URL_PREFIX}/`;
export const registerWithGroupPathWeb = `${WEBAPP_URL_PREFIX}/registerWithGroup/*`;

export function getLandingPath(isAdminApp) {
  return isAdminApp ? landingPath : landingPathWeb;
}

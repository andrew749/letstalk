const NOTIF_OPENED_PREFIX = 'notif-opened';
const LINK_CLICKED_PREFIX = 'link-clicked';

/**
 * Initializes mixpanel metrics, linking the metrics to the given userId.
 *
 */
function identifyUser(userId) {
  mixpanel.identify(userId);
}

function getProperties(extraProperties) {
  let properties = window.passedContext;
  properties = Object.assign(properties, extraProperties || {});
  properties.referrer = document.referrer;
  return properties;
}

function trackNotificationOpened(notificationName, extraProperties) {
  const eventName = notificationName + '/' + NOTIF_OPENED_PREFIX;
  mixpanel.track(eventName, getProperties(extraProperties));
}

function trackAction(notificationName, action, extraProperties) {
  const eventName = notificationName + '/' + action;
  mixpanel.track(eventName, getProperties(extraProperties));
}

function trackLinkClicked(notificationName, linkAccessor, linkName, extraProperties) {
  const eventName = notificationName + '/' + LINK_CLICKED_PREFIX + '/' + linkName;
  mixpanel.track_links(linkAccessor, eventName, getProperties(extraProperties));
}

export { identifyUser, trackNotificationOpened, trackAction, trackLinkClicked };

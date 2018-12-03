const NOTIF_OPENED_PREFIX = 'notif-opened';
const LINK_CLICKED_PREFIX = 'link-clicked';

/**
 * Initializes mixpanel metrics, linking the metrics to the given userId.
 *
 */
function identifyUser(userId) {
  mixpanel.identify(userId);
}

function trackNotificationOpened(notificationName, properties) {
  const eventName = notificationName + '/' + NOTIF_OPENED_PREFIX;
  properties = properties || window.passedContext;
  mixpanel.track(eventName, properties);
}

function trackLinkClicked(notificationName, linkAccessor, linkName, properties) {
  const eventName = notificationName + '/' + LINK_CLICKED_PREFIX + '/' + linkName;
  properties = properties || window.passedContext;
  properties.referrer = document.referrer;
  mixpanel.track_links(linkAccessor, eventName, properties);
}

export { identifyUser, trackNotificationOpened, trackLinkClicked };

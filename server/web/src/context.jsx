/**
 * Context that is passed to the react app from outside React-land.
 */
export default function getProperty(propertyKey) {
  if (window.passedContext.hasOwnProperty(propertyKey)) {
    return window.passedContext[propertyKey];
  }
  console.warn(`Property ${propertyKey} not set`)
  return null;
}

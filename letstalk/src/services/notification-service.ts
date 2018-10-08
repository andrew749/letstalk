import { Store } from 'redux';
import {
  NavigationActions,
  NavigationContainerComponent,
  NavigationParams,
} from 'react-navigation';
import { Linking } from 'expo';

import { RootState } from '../redux';
import { fetchBootstrap } from '../redux/bootstrap/reducer';
import { Notification } from '../models/notification';

export interface RawNotification {
  readonly data: Notification;
  readonly origin: 'selected' | 'received';
}

export default class NotificationService {
  private notification: any;
  private store: Store<RootState>;
  private navContainer: NavigationContainerComponent;

  constructor(store: Store<RootState>) {
    this.notification = null;
    this.store = store;
    this.navContainer = null;

    this.onPress = this.onPress.bind(this);
  }

  setNotifContainer(notifContainer: any) {
    this.notification = notifContainer;
  }

  setNavContainer(navContainer: NavigationContainerComponent) {
    this.navContainer = navContainer;
  }

  navigate(routeName: string, params?: NavigationParams) {
    if (!!this.navContainer) {
      this.navContainer.dispatch(NavigationActions.navigate({ routeName, params }));
    }
  }

  navigateToLink(link: string) {
    if (!!link) {
      let { path, queryParams } = Linking.parse(link);
      this.navContainer.dispatch(NavigationActions.navigate({
        routeName: path,
        params: queryParams,
      }));
    }
  }

  private async actOnNotification(notification: Notification): Promise<void> {
    switch (notification.type) {
      case 'ADHOC_NOTIFICATION':
      case 'NEW_MATCH':
      case 'NEW_CREDENTIAL_MATCH':
      case 'CONNECTION_REQUESTED':
      case 'CONNECTION_ACCEPTED':
        break;
      default:
        // Ensure exhaustiveness of select
        const _: never = notification;
        // This case could happen, but we wouldn't do anything anyways
    }
  }

  private onPress(notification: Notification) {
    // update action to use deeplink
    return () => {
      switch (notification.type) {
        case 'ADHOC_NOTIFICATION':
        case 'NEW_MATCH':
        case 'NEW_CREDENTIAL_MATCH':
        case 'CONNECTION_REQUESTED':
        case 'CONNECTION_ACCEPTED':
          this.navigateToLink(notification.link);
          break;
        default:
          // Ensure exhaustiveness of select
          const _: never = notification;
          // This case could happen, but we wouldn't do anything anyways
      }
    };
  }

  async handleNotification(notification: RawNotification): Promise<void> {
    if (!!this.notification) {
      if (notification.origin === 'received') {
        this.notification.show({
          title: notification.data.title,
          message: notification.data.message,
          onPress: this.onPress.bind(this, notification.data),
        });
      } else if (notification.origin === 'selected') {
        this.onPress(notification.data);
      }
    }
    await this.actOnNotification(notification.data);
  }
}

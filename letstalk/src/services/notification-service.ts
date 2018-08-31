import { Store } from 'redux';
import {
  NavigationActions,
  NavigationContainerComponent,
  NavigationParams,
} from 'react-navigation';

import { RootState } from '../redux';
import { credentialRequestRemove } from '../redux/credential-requests/actions';
import { credentialRemove } from '../redux/credentials/actions';
import { fetchBootstrap } from '../redux/bootstrap/reducer';

interface BaseNotificationData {
  readonly title: string;
  readonly message: string;
}

const SIDE_ASKER = "ASKER";
const SIDE_ANSWERER = "ANSWERER";

interface RequestToMatchNotificationData extends BaseNotificationData {
  readonly type: 'REQUEST_TO_MATCH';
  readonly side: string;
  readonly requestId: number;
}

interface AdHocNotification extends BaseNotificationData {
  readonly type: 'ADHOC_NOTIFICATION';
}

interface NewMatchNotification extends BaseNotificationData {
  readonly type: 'NEW_MATCH';
}

interface NewCredentialMatchNotification extends BaseNotificationData {
  readonly type: 'NEW_CREDENTIAL_MATCH';
}

interface NewMatchVerifiedNotificaiton extends BaseNotificationData {
  readonly type: 'MATCH_VERIFIED';
}

type NotificationData = RequestToMatchNotificationData
| AdHocNotification
| NewMatchNotification
| NewCredentialMatchNotification
| NewMatchVerifiedNotificaiton;

export interface Notification {
  readonly data: NotificationData;
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

  private async actOnNotification(notification: Notification): Promise<void> {
    const data = notification.data;
    switch (data.type) {
      case 'REQUEST_TO_MATCH':
        if (data.side === SIDE_ASKER) {
          await this.store.dispatch(credentialRequestRemove(data.requestId));
        }
        await fetchBootstrap()(this.store.dispatch, null, null);
        break;
      case 'ADHOC_NOTIFICATION':
      case 'NEW_MATCH':
      case 'NEW_CREDENTIAL_MATCH':
      case 'MATCH_VERIFIED':
      default:
        // Ensure exhaustiveness of select
        // const _: never = data.type;
        // This case could happen, but we wouldn't do anything anyways
    }
  }

  private onPress(notification: Notification) {
    return () => {
      switch (notification.data.type) {
        case 'REQUEST_TO_MATCH':
          // TODO: Make this do a reset instead
          this.navigate('Home');
          break;
        case 'ADHOC_NOTIFICATION':
        case 'NEW_MATCH':
        case 'NEW_CREDENTIAL_MATCH':
        case 'MATCH_VERIFIED':
          break;
        default:
          // Ensure exhaustiveness of select
          // const _: never = notification.data.Type;
          // This case could happen, but we wouldn't do anything anyways
      }
    };
  }

  async handleNotification(notification: Notification): Promise<void> {
    if (!!this.notification) {
      this.notification.show({
        title: notification.data.title,
        message: notification.data.message,
        onPress: this.onPress(notification),
      });
    }

    await this.actOnNotification(notification);
  }
}

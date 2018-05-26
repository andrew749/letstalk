import { Store } from 'redux';

import { RootState } from '../redux';
import { credentialRequestRemove } from '../redux/credential-requests/actions';
import { credentialRemove } from '../redux/credentials/actions';

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

type NotificationData = RequestToMatchNotificationData;

export interface Notification {
  readonly data: NotificationData;
}

export default class NotificationService {
  private notification: any
  private store: Store<RootState>

  constructor(notification: any, store: Store<RootState>) {
    this.notification = notification;
    this.store = store;
  }

  private async actOnNotification(notification: Notification): Promise<void> {
    const data = notification.data;
    console.log(data);
    switch (data.type) {
      case 'REQUEST_TO_MATCH':
        switch (data.side) {
          case SIDE_ASKER:
            await this.store.dispatch(credentialRequestRemove(data.requestId));
            break;
          case SIDE_ANSWERER:
            await this.store.dispatch(credentialRemove(data.requestId));
            break;
        }
        break;
      default:
        // Ensure exhaustiveness of select
        const _: never = data.type;
        // This case could happen, but we wouldn't do anything anyways
    }
  }

  async handleNotification(notification: Notification): Promise<void> {
    this.notification.show({
      title: notification.data.title,
      message: notification.data.message,
    });

    await this.actOnNotification(notification);
  }
}


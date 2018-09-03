export type NotifState =
  | 'UNREAD'
  | 'READ';

interface BaseNotification {
  notificationId: number;
  userId: number;
  type: string;
  state: NotifState;
  title: string;
  message: string;
  timestamp: Date;
  thumbnail: string;
}

interface NewCredentialMatchData {
  credentialName: string,
  userName: string,
  side: "ASKER" | "ANSWERER",
}

export interface NewCredentialMatchNotification extends BaseNotification {
  type: 'NEW_CREDENTIAL_MATCH',
  data: NewCredentialMatchData,
}

interface AdhocNotificationData {

}

export interface AdhocNotification extends BaseNotification {
  type: 'ADHOC_NOTIFICATION',
  data: AdhocNotificationData,
}

export interface NewMatchData {

}

export interface NewMatchNotification extends BaseNotification {
  type: 'NEW_MATCH',
  data: NewMatchData,
}

export type Notification =
  | NewCredentialMatchNotification
  | AdhocNotification
  | NewMatchNotification;

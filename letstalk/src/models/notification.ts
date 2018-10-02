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
  link: string;
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

export interface ConnectionRequestedData {
}

export interface ConnectionRequestedNotification extends BaseNotification {
  type: 'CONNECTION_REQUESTED',
  data: ConnectionRequestedData,
}

export interface ConnectionAcceptedData {
}

export interface ConnectionAcceptedNotification extends BaseNotification {
  type: 'CONNECTION_ACCEPTED',
  data: ConnectionAcceptedData,
}

export type Notification =
  | NewCredentialMatchNotification
  | AdhocNotification
  | NewMatchNotification
  | ConnectionRequestedNotification
  | ConnectionAcceptedNotification;

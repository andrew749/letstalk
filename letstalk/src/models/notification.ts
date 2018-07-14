export type NotifState =
  | 'UNREAD'
  | 'READ';

interface BaseNotification {
  notificationId: number;
  userId: number;
  state: NotifState;
  createdAt: Date;
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

export type Notification =
  | NewCredentialMatchNotification;

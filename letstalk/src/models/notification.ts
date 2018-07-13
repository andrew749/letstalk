export interface Notification {
  notificationId: number;
  userId: number;
  type: string;
  state: string;
  data: object;
  createdAt: Date;
}

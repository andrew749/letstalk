// TODO: Figure out how to make this library work, since we would prefer to encrypt our token
// import * as Keychain from 'react-native-keychain';
import { Alert, AsyncStorage, Linking, Platform } from 'react-native';

import requestor from './requests';
import { SessionService, SessionToken, RemoteSessionService } from './session-service';
import {FORGOT_PASSWORD_ROUTE, SEND_EMAIL_VERIFICATION_ROUTE} from './constants';
import {Notifications, Permissions} from "expo";
import {SendAccountVerificationEmailRequest} from "../models/verify_email";

export class Auth {
  private sessionService: SessionService
  private sessionToken: SessionToken

  private static readonly tokenLocation = '@app/session';

  constructor(sessionService: SessionService) {
    this.sessionService = sessionService;
    this.sessionToken = null;
  }

  async getSessionToken(): Promise<SessionToken> {
    if (this.sessionToken === null) {
      this.sessionToken = await AsyncStorage.getItem(Auth.tokenLocation);
    }
    return this.sessionToken;
  }

  private async setSessionToken(sessionToken: SessionToken): Promise<void> {
    await AsyncStorage.setItem(Auth.tokenLocation, sessionToken);
    this.sessionToken = sessionToken;
  }

  async login(email: string, password: string, notificationToken?: string): Promise<void> {
    const sessionToken = await this.sessionService.login(
        email,
        password,
        notificationToken,
    );
    await this.setSessionToken(sessionToken);
  }

  // Returns whether the login completed.
  async loginWithFb(notificationToken?: string): Promise<boolean> {
    const sessionToken = await this.sessionService.loginWithFb(notificationToken);
    if (sessionToken === null) return false;
    await this.setSessionToken(sessionToken);
    return true;
  }

  async logout(): Promise<void> {
    const sessionToken = await this.getSessionToken();
    AsyncStorage.removeItem(Auth.tokenLocation);
    await this.sessionService.logout(sessionToken);
  }

  async forgotPassword(email: string): Promise<any> {
    return await requestor.post(FORGOT_PASSWORD_ROUTE, {"email": email});
  }

  async sendVerificationEmail(email: string): Promise<void> {
    const sessionToken = await auth.getSessionToken();
    const req: SendAccountVerificationEmailRequest = {"email": email};
    const resp = await requestor.post(SEND_EMAIL_VERIFICATION_ROUTE, req, sessionToken);
  }

  async linkFB(): Promise<boolean> {
    const sessionToken = await this.getSessionToken();
    return await this.sessionService.linkFb(sessionToken);
  }

  // showModal: if true, and user on ios is missing the notification permission, will show a modal
  // asking them to change the settings.
  async registerForPushNotificationsAsync(showModalOnMissing: boolean = false): Promise<string> {
    const { status: existingStatus } = await Permissions.getAsync(Permissions.NOTIFICATIONS as any);
    let finalStatus = existingStatus;

    // only ask if permissions have not already been determined, because
    // iOS won't necessarily prompt the user a second time.
    if (existingStatus !== 'granted') {
      if (Platform.OS === 'ios' && showModalOnMissing) {
        const onPress = () => Linking.openURL('app-settings:');
        // iOS doesn't show dialog a second time, so refer users to app settings to change config.
        Alert.alert(
          'Notification Permissions',
          'Open app settings to enable notifications permissions. Notifications will let you know ' +
          'when you have new matches, remind you to check-in with your connections and inform you ' +
          'about events happening on campus.',
          [
            {text: 'Cancel', onPress: () => null, style: 'cancel'},
            {text: 'Open Settings', onPress: onPress, style: 'default'},
          ],
        );
      }
      // Android remote notification permissions are granted during the app
      // install, so this will only ask on iOS
      const { status } = await Permissions.askAsync(Permissions.NOTIFICATIONS as any);
      finalStatus = status;
    }

    // Stop here if the user did not grant permissions
    if (finalStatus !== 'granted') {
      return null;
    }

    // Get the token that uniquely identifies this device
    let token = await Notifications.getExpoPushTokenAsync();
    console.log("Registered with expo notification service: " + token);
    return token;
  }
};

const auth = new Auth(new RemoteSessionService(requestor));

export default auth;

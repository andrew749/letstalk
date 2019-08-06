import { Alert, Linking, Platform } from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import * as Permissions from 'expo-permissions';
import requestor, { Requestor } from './requests';
import auth, { Auth } from './auth';
import { PROFILE_PIC_UPLOAD_ROUTE } from './constants';

interface PhotoService {
  uploadProfilePhoto(uri: string): void;
  getPhotoFromPicker(): Promise<PhotoResult>;
}

export type PhotoResult = {
  uri: string,
}

export class PhotoServiceImpl implements PhotoService {
  private _requestor: Requestor;
  constructor(requestor: Requestor) {
    this._requestor = requestor;
  }

  async getPhotoFromPicker(): Promise<PhotoResult> {
    // Display the camera to the user and wait for them to take a photo or to cancel
    // the action
    let cameraPermission = await Permissions.askAsync(Permissions.CAMERA_ROLL as any);
    if (cameraPermission.status !== 'granted') {
      if (Platform.OS === 'ios') {
        const onPress = () => Linking.openURL('app-settings:');
        // iOS doesn't show dialog a second time, so refer users to app settings to change config.
        Alert.alert(
          'Camera Roll Permissions',
          'Open app settings to enable camera roll permissions',
          [
            {text: 'Cancel', onPress: () => null, style: 'cancel'},
            {text: 'Open Settings', onPress: onPress, style: 'default'},
          ],
        );
        cameraPermission = await Permissions.askAsync(Permissions.CAMERA_ROLL as any);
      }
      if (cameraPermission.status !== 'granted') return;
    }

    let result = await ImagePicker.launchImageLibraryAsync({
      allowsEditing: true,
      aspect: [1, 1],
      quality: 0.2,
      base64: true,
    });

    if (result.cancelled === true) {
      return;
    }

    return {
      uri: result.uri,
    };
  }

  async uploadProfilePhoto(uri: string) {
    // Upload the image using the fetch and FormData APIs
    let formData = new FormData();
    const sessionToken = await auth.getSessionToken();
    const data = await fetch(uri);

    // sends a base 64 encoded string
    formData.append('photo', await data.blob());

    return await this._requestor.postFormData(PROFILE_PIC_UPLOAD_ROUTE, formData, sessionToken);
  }
}

export const photoService = new PhotoServiceImpl(requestor);
export default photoService;

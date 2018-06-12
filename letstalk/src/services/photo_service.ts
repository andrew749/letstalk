import {ImagePicker, FileSystem, Permissions} from 'expo';
import requestor, {Requestor} from './requests';
import auth, {Auth} from './auth';
import {PROFILE_PIC_UPLOAD_ROUTE} from './constants';

interface PhotoService {
  uploadProfilePhoto(uri: string): void;
  getPhotoFromPicker(): Promise<PhotoResult>;
}

export type PhotoResult = {
  uri: string,
  data: string, // base64 encoded
}

export class PhotoServiceImpl implements PhotoService {
  private _requestor: Requestor;
  constructor(requestor: Requestor) {
    this._requestor = requestor;
  }

  async getPhotoFromPicker(): Promise<PhotoResult> {
    // Display the camera to the user and wait for them to take a photo or to cancel
    // the action
    const cameraPermission = await Permissions.askAsync(Permissions.CAMERA);

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
      data: result.base64,
    };
  }

  async uploadProfilePhoto(uri: string) {
    // Upload the image using the fetch and FormData APIs
    let formData = new FormData();
    let photoData = await FileSystem.readAsStringAsync(uri);

    // sends a base 64 encoded string
    formData.append('photo', photoData);

    return await this._requestor.postBinary(PROFILE_PIC_UPLOAD_ROUTE, {
      body: formData,
      headers: {
        'content-type': 'multipart/form-data',
      },
    }, null);
  }
}

export const photoService = new PhotoServiceImpl(requestor);
export default photoService;

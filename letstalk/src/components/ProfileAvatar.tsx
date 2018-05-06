import React from 'react';
import { Avatar, AvatarProps, FormInput, FormInputProps } from 'react-native-elements';
import { ImageURISource, StyleSheet } from 'react-native';
import photoService, {PhotoResult} from '../services/photo_service';
import { WrappedFieldProps } from 'redux-form';
import {
    FormProps
} from '../components'

interface ProfileAvatarProps extends AvatarProps {
  userId?: string,
  overrideUri?: string,
  editable: boolean,
}

type Props = WrappedFieldProps & ProfileAvatarProps;

function getProfilePicUrl(userId: string): string {
  return `https://s3.amazonaws.com/hive-user-profile-pictures/{userId}`;
}

type State = {
  uri: string
}

/**
 * A component to be used in the app to render profile pictures for a
 * user.
 *
 * If there isn't a profile pic then fallback
 */
class ProfileAvatar extends React.Component<Props, State> {
  render() {
    let avatarSource = undefined;
    let props = this.props;
    let overrideUri = props.overrideUri;
    if (props.input.value) {
      overrideUri = (props.input.value as PhotoResult).uri;
    }
    let onChange = this.props.input.onChange;

    if (overrideUri) {
      const profilePicUrl = overrideUri;
      avatarSource = {uri: profilePicUrl};
    } else if (props.userId) {
      const profilePicUrl = getProfilePicUrl(props.userId);
      avatarSource = {uri: profilePicUrl};
    }

    let pressAction = () => {};
    if (props.editable) {
      pressAction = async() => {
        let photoResult = await photoService.getPhotoFromPicker();
        onChange(photoResult);
      };
    }

    return (
      <Avatar
        {...props}
        large
        rounded
        // default
        icon={{name: 'person'}}
        source={ avatarSource }
        onPress={ pressAction }
        activeOpacity={0.7}
      />
    );
  }
}

const styles = StyleSheet.create({

});

export default ProfileAvatar;

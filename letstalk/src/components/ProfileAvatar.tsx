import React from 'react';
import { Avatar, AvatarProps, FormInput, FormInputProps } from 'react-native-elements';
import { ImageURISource, StyleSheet } from 'react-native';
import photoService, {PhotoResult} from '../services/photo_service';
import {profileService, RemoteProfileService} from '../services/profile-service';
import { WrappedFieldProps } from 'redux-form';
import {
    FormProps
} from '../components'

import Colors from '../services/colors';

interface ProfileAvatarProps extends AvatarProps {
  userId?: string,
}

/**
 * A component to be used in the app to render profile pictures for a
 * user.
 *
 * If there isn't a profile pic then fallback
 */

interface ProfileAvatarState {
  avatarSource: ImageURISource,
}

class ProfileAvatar extends React.Component<ProfileAvatarProps, ProfileAvatarState> {

  constructor(props: ProfileAvatarProps) {
    super(props);
    this.state={
      avatarSource: props.source
    }
  }

  async getDerivedStateFromProps(props: ProfileAvatarProps) {
    if (props.source) {
      this.setState({avatarSource: props.source});
    }
  }

  async componentDidMount() {
    let props = this.props;
    if (props.userId) {
      const profilePicUrl = await profileService.getProfilePicUrl(props.userId);
      if (profilePicUrl != undefined || profilePicUrl != null) {
        this.setState({ avatarSource: {uri: profilePicUrl} });
      }
    }
  }

  render() {
    let props = this.props;
    return (
      <Avatar
        {...props}
        xlarge
        rounded
        // default
        icon={{name: 'person'}}
        source={ this.state.avatarSource }
        activeOpacity={0.7}
      />
    );
  }
}

type FormElementProps = WrappedFieldProps & ProfileAvatarProps;
export class ProfileAvatarEditableFormElement extends React.Component<FormElementProps> {
  render() {
    let props = this.props;
    let avatarSource = undefined;

    // handle an on click
    let onChange = this.props.input.onChange;

    let pressAction = async() => {
        let photoResult = await photoService.getPhotoFromPicker();
        onChange(photoResult);
      };

    // the user changed the form contents to an image
    if (props.input.value) {
      let uri = (props.input.value as PhotoResult).uri;
      avatarSource = {uri: uri};
    }
    return (
        <ProfileAvatar
          {...props}
          onPress={ pressAction }
          source={ avatarSource }
        />
    );
  }
}

const styles = StyleSheet.create({

});

export default ProfileAvatar;

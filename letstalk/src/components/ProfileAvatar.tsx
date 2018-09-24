import React from 'react';
import { Avatar, AvatarProps, FormInput, FormInputProps } from 'react-native-elements';
import { ImageURISource, StyleProp, StyleSheet, View, ViewStyle } from 'react-native';
import photoService, {PhotoResult} from '../services/photo_service';
import {profileService, RemoteProfileService} from '../services/profile-service';
import { WrappedFieldProps } from 'redux-form';
import {
    FormProps
} from '../components'
import { MaterialIcons, MaterialCommunityIcons } from '@expo/vector-icons';

import Colors from '../services/colors';

interface ProfileAvatarProps extends AvatarProps {
  userId?: string,
  edit?: boolean
}

interface ProfileAvatarFormProps {
  onChange: any;
  uri: any;
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

  async componentDidUpdate(prevProps: ProfileAvatarProps, prevState: ProfileAvatarState) {
    if (prevProps.source !== this.props.source) {
      this.setState({avatarSource: this.props.source});
    }
  }

  async componentDidMount() {
    let props = this.props;
    if (props.userId) {
      const profilePicUrl = await profileService.getProfilePicUrl(props.userId);
      if (!!profilePicUrl) {
        this.setState({ avatarSource: {uri: profilePicUrl} });
      }
    }
  }

  render() {
    let props = this.props;
    return (
      <View>
        <Avatar
          {...props}
          xlarge
          rounded
          // default
          icon={{name: 'person'}}
          source={ this.state.avatarSource }
          activeOpacity={0.7}
        />
        {props.edit && <MaterialIcons
          style={styles.editAvatar}
          name="camera-alt"
          size={25}
          color={Colors.WHITE}
          onPress={this.props.onPress}
        />}
      </View>
    );
  }
}

type FormElementProps = WrappedFieldProps & ProfileAvatarProps & ProfileAvatarFormProps;
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
    if (props.input.value && props.input.value.uri) {
      let uri = (props.input.value as PhotoResult).uri;
      avatarSource = {uri: uri};
    } else {
      avatarSource = props.uri;
    }
    return (
        <ProfileAvatar
          {...props}
          xlarge
          edit
          onPress={ pressAction }
          source={ avatarSource }
          containerStyle= {styles.profilePicture}
        />
    );
  }
}

const styles = StyleSheet.create({
  editAvatar: {
    position: 'absolute',
    right: 23,
    bottom: 23,
    padding: 5,
    backgroundColor: Colors.HIVE_SUBDUED,
    borderRadius: 30,
  },
  profilePicture: {
    margin: 20
  },
});

export default ProfileAvatar;

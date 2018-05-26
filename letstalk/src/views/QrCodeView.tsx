import React, { Component } from 'react';
import { View } from 'react-native';
import {ActionCreator, connect} from "react-redux";
import {ThunkAction} from "redux-thunk";
import {ActionTypes} from "../redux/profile/actions";
import { State as ProfileState, fetchProfile } from '../redux/profile/reducer';
import {NavigationScreenProp, NavigationStackAction} from "react-navigation";
import {RootState} from "../redux/index";
import QRCode from 'react-native-qrcode';

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
}

interface Props extends ProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class QrCodeView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'My Code',
  });

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
  }

  async componentDidMount() {
    await this.load();
  }

  private async load() {
    await this.props.fetchProfile();
  }

  render() {
    const { secret } = this.props.profile;
    return (
      <View>
        {!!secret && <QRCode
          value={secret}
          size={200}
          bgColor='black'
          fgColor='white'/>
        }
      </View>
    );
  }
}

export default connect(({profile}: RootState) => profile, { fetchProfile })(QrCodeView);

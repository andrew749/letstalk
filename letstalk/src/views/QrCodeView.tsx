import React, { Component } from 'react';
import {
  Text,
  View,
} from 'react-native';
import {ActionCreator, connect} from "react-redux";
import {ThunkAction} from "redux-thunk";
import {ActionTypes} from "../redux/bootstrap/actions";
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import {NavigationScreenProp, NavigationStackAction} from "react-navigation";
import {RootState} from "../redux/index";
import QRCode from 'react-native-qrcode';

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class QrCodeView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'QrCode',
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
  }

  private async load() {
    await this.props.fetchBootstrap();
  }

  render() {
    const { secret } = this.props.bootstrap && this.props.bootstrap.me;
    return (
      <View>
        {secret && <QRCode
          value={secret}
          size={200}
          bgColor='black'
          fgColor='white'/>
        }
      </View>
    );
  }
}

export default connect(({bootstrap}: RootState) => bootstrap, { fetchBootstrap })(QrCodeView);

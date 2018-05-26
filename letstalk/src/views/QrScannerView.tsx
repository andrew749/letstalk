import React, { Component } from 'react';
import {
  Alert,
  Linking,
  Dimensions,
  LayoutAnimation,
  Text,
  View,
  StatusBar,
  StyleSheet,
  TouchableOpacity,
} from 'react-native';
import {BarCodeScanner, Permissions} from 'expo';
import { ToastActionsCreators } from 'react-native-redux-toast';
import { errorToast, infoToast } from '../redux/toast';
import {NavigationScreenProp, NavigationStackAction} from "react-navigation";
import {ActionTypes} from "../redux/bootstrap/actions";
import {ThunkAction} from "redux-thunk";
import {ActionCreator, connect, Dispatch} from "react-redux";
import {RootState} from "../redux/index";
import { State as BootstrapState, fetchBootstrap } from "../redux/bootstrap/reducer";

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class QrScannerView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'QrScanner',
  })

  state = {
    hasCameraPermission: false,
    lastScannedUrl: '',
  };

  constructor(props: Props) {
    super(props);
    console.log("QrScanner constructor");

    this.load = this.load.bind(this);
  }

  async componentDidMount() {
    this.load();
  }

  private async load() {
    await Promise.all([
      this._requestCameraPermission(),
      this.props.fetchBootstrap(),
    ]);
  }

  _requestCameraPermission = async () => {
    const { status } = await Permissions.askAsync(Permissions.CAMERA);
    console.log("permissions result: ", status);
    this.setState({
      hasCameraPermission: status === 'granted',
    });
  };

  _handleBarCodeRead = async (result: { type: string; data: string; }) => {
    if (result.data !== this.state.lastScannedUrl) {
      LayoutAnimation.spring();
      this.setState({ lastScannedUrl: result.data });
      await this.props.infoToast(`Scanned barcode ${result.data}`);
    }
  };

  render() {
    return (
      <View style={styles.container}>

        {!this.state.hasCameraPermission
          ? <Text style={{ color: '#fff' }}>
            Camera permission required.
          </Text>
          // Need to use React.createElement to skip type checking on style prop.
          : React.createElement(BarCodeScanner, {
            onBarCodeRead: this._handleBarCodeRead,
            // @ts-ignore bug where style isn't detected in BarCodeScannerProps
            style: {
              height: Dimensions.get('window').height,
              width: Dimensions.get('window').width,
            },
          })
        }

        <View style={styles.bottomBar}>
          <Text numberOfLines={1} style={styles.bottomText}>
            {'Scan a QR code to confirm a match or meeting'}
          </Text>
        </View>

        <StatusBar hidden />
      </View>
    );
  };
}

export default connect(({bootstrap}: RootState) => bootstrap, { fetchBootstrap, errorToast, infoToast })(QrScannerView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#000',
  },
  bottomBar: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(0,0,0,0.5)',
    padding: 15,
    flexDirection: 'row',
  },
  bottomText: {
    color: '#fff',
    fontSize: 14,
  },
});

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
import {NavigationActions, NavigationScreenProp, NavigationStackAction} from "react-navigation";
import {ActionTypes} from "../redux/bootstrap/actions";
import {ThunkAction} from "redux-thunk";
import {ActionCreator, connect, Dispatch} from "react-redux";
import {RootState} from "../redux/index";
import { State as BootstrapState, fetchBootstrap } from "../redux/bootstrap/reducer";
import meetingService from "../services/meeting";

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
    lastScannedBarcode: '',
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
      this.requestCameraPermission(),
      this.props.fetchBootstrap(),
    ]);
  }

  requestCameraPermission = async () => {
    const { status } = await Permissions.askAsync(Permissions.CAMERA);
    console.log("permissions result: ", status);
    this.setState({
      hasCameraPermission: status === 'granted',
    });
  };

  handleBarCodeRead = async (result: { type: string; data: string; }) => {
    const barcode = result.data;
    if (barcode !== this.state.lastScannedBarcode) {
      LayoutAnimation.spring();
      await this.setState({ lastScannedBarcode: barcode });
      try {
        await meetingService.postMeetingConfirmation({secret: barcode});
        await this.props.infoToast('Meeting confirmed!');
        await this.props.navigation.dispatch(NavigationActions.back());
      } catch(error) {
        await this.props.errorToast('Failed to confirm meeting, please try again.')
        await this.setState({ lastScannedBarcode: '' });
      }
    }
  };

  render() {
    return (
      <View style={styles.container}>

        {!this.state.hasCameraPermission
          ? <Text style={{color: '#fff'}}>
            Camera permission required.
          </Text>
          // Need to use React.createElement to skip type checking on style prop.
          : React.createElement(BarCodeScanner, {
            onBarCodeRead: this.handleBarCodeRead,
            // @ts-ignore bug where style isn't detected in BarCodeScannerProps
            style: {
              height: Dimensions.get('window').height,
              width: Dimensions.get('window').width,
            },
          })
        }

        <View style={styles.bottomBar}>
          <Text numberOfLines={1} style={styles.bottomText}>
            {'Scan a QR code to confirm a meeting'}
          </Text>
        </View>

        <StatusBar hidden/>
      </View>
    );
  }
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
    fontSize: 18,
  },
});

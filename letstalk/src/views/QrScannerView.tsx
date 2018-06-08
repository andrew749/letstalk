import React, { Component } from 'react';
import {
  Dimensions,
  LayoutAnimation,
  Text,
  View,
  StatusBar,
  StyleSheet,
} from 'react-native';
import {BarCodeScanner, Permissions} from 'expo';
import { ToastActionsCreators } from 'react-native-redux-toast';
import { errorToast, infoToast } from '../redux/toast';
import {NavigationActions, NavigationScreenProp, NavigationStackAction} from "react-navigation";
import {connect, Dispatch} from "react-redux";
import {RootState} from "../redux/index";
import meetingService from "../services/meeting";
import { headerStyle } from './TopHeader';
import { AnalyticsHelper } from '../services';

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class QrScannerView extends Component<Props> {
  QR_SCANNER_VIEW_IDENTIFIER = "QrScannerView";

  static navigationOptions = () => ({
    headerTitle: 'Scan A Code',
    headerStyle,
  })

  state = {
    hasCameraPermission: false,
    lastScannedBarcode: '',
  };

  constructor(props: Props) {
    super(props);
    this.load = this.load.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.QR_SCANNER_VIEW_IDENTIFIER);
    await this.load();
  }

  private async load() {
    await this.requestCameraPermission();
  }

  requestCameraPermission = async () => {
    const { status } = await Permissions.askAsync(Permissions.CAMERA);
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
        await this.props.errorToast('Failed to confirm meeting, please try again')
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

export default connect(null, { errorToast, infoToast })(QrScannerView);

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

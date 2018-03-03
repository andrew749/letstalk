// NB: This is meant for iOS only.
// TODO: Rename to .ios.tsx
import React, { ReactNode } from 'react';
import {
  Dimensions,
  Modal,
  DatePickerIOS,
  DatePickerIOSProperties,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';
import Card from './Card';

const SCREEN_WIDTH = Dimensions.get('window').width;

const styles = StyleSheet.create({
  basicContainer: {
    flex: 1,
    justifyContent: 'flex-end',
    alignItems: 'center',
  },

  overlayContainer: {
    flex: 1,
    width: SCREEN_WIDTH,
  },

  modalContainer: {
    width: SCREEN_WIDTH,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 0,
    backgroundColor: '#F5FCFF',
  },

  buttonView: {
    width: SCREEN_WIDTH,
    padding: 8,
    borderTopWidth: 0.5,
    borderTopColor: 'lightgrey',
    justifyContent: 'space-between',
    flexDirection: 'row',
  },

  label: {
    fontWeight: 'bold',
    fontSize: 18,
    textAlign: 'center',
  },

  valueLabel: {
    fontSize: 18,
    textAlign: 'center',
    color: 'gray',
    paddingLeft: 10,
  },

  cardWithValue: {
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'flex-start',
    height: 60,
    padding: 10,
    paddingLeft: 20,
  },

  cardWithoutValue: {
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    height: 60,
    padding: 10,
  },
});

type Props = WrappedFieldProps & {
  label: string;
  valueLabel?: string;
  children: ReactNode;
  onSubmitPress(): void;
};

interface State {
  modalVisible: boolean;
};

class StatefulBottomModal extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {
      modalVisible: false,
    };

    this.show = this.show.bind(this);
    this.hide = this.hide.bind(this);
    this.hideAndBlur = this.hideAndBlur.bind(this);
    this.onSubmitPress = this.onSubmitPress.bind(this);
  }

  onSubmitPress() {
    this.props.onSubmitPress();
    (this.props.input.onBlur as () => void)();
    this.hide();
  }

  hideAndBlur() {
    // So that we show required error if user clicks 'Cancel'
    (this.props.input.onBlur as () => void)();
    this.hide();
  }

  show() {
    this.setState({
      modalVisible: true,
    });
  }

  hide() {
    this.setState({
      modalVisible: false,
    });
  }

  renderDisplayWithValue() {
    const { valueLabel, label } = this.props;
    return (
      <View style={{flex: 1, flexDirection: 'row', alignItems: 'center'}}>
        <Text style={styles.label}>{label}</Text>
        <Text style={styles.valueLabel}>{valueLabel}</Text>
      </View>
    );
  }

  renderDisplayWithoutValue() {
    const { label } = this.props;
    return (
        <Text style={styles.label}>{label}</Text>
    );
  }

  render() {
    const { children } = this.props;
    const { value } = this.props.input;
    const { error, touched, warning } = this.props.meta;
    const { modalVisible } = this.state;

    const display = value ? this.renderDisplayWithValue() : this.renderDisplayWithoutValue();
    const cardStyle = value ? styles.cardWithValue : styles.cardWithoutValue;

    // TODO: Add onPress to Card
    return (
      <View>
        <Modal
          animationType={'slide'}
          transparent
          visible={modalVisible}
        >
          <View style={styles.basicContainer}>
            <View style={styles.overlayContainer}>
              <TouchableWithoutFeedback onPress={this.hideAndBlur}>
                <View style={styles.overlayContainer} />
              </TouchableWithoutFeedback>
            </View>
            <View style={styles.modalContainer}>
              <View style={styles.buttonView}>
                <TouchableOpacity onPress={this.hideAndBlur}>
                  <Text>Cancel</Text>
                </TouchableOpacity>
                <TouchableOpacity onPress={this.onSubmitPress}>
                  <Text>Confirm</Text>
                </TouchableOpacity>
              </View>
              <View>
                {children}
              </View>
            </View>
          </View>
        </Modal>
        <TouchableOpacity onPress={this.show}>
          <Card style={cardStyle}>
              {display}
              {touched && (
                (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
                (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
          </Card>
        </TouchableOpacity>
      </View>
    );
  }
}

const BottomModal: React.SFC<Props> = props => {
  return <StatefulBottomModal {...props} />;
}
export default BottomModal;

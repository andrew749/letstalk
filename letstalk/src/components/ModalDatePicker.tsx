// TODO: Remove shared code between this and [ModalPicker]
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
import ActionButton from './ActionButton';

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

  bottomPicker: {
    width: SCREEN_WIDTH,
  },
});

type Props = WrappedFieldProps & {
  label: string;
  defaultDate: Date;
  mode?: 'date' | 'time' | 'datetime';
};

interface State {
  modalVisible: boolean;
};

class StatefulModalDatePicker extends React.Component<Props, State> {
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
    const { value, onChange } = this.props.input;
    onChange(value || this.props.defaultDate);
    this.hide()
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

  render() {
    const { defaultDate, label, mode } = this.props;
    const { onChange, value } = this.props.input;
    const { error, touched, warning } = this.props.meta;
    const { modalVisible } = this.state;

    // TODO: Don't use action button
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
                <DatePickerIOS
                  mode={mode}
                  style={styles.bottomPicker}
                  date={value || defaultDate}
                  onDateChange={onChange}
                >
                </DatePickerIOS>
              </View>
            </View>
          </View>
        </Modal>
        <ActionButton
          title={value ? value.toString() : label}
          onPress={this.show}
        />
        {touched && (
          (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
          (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
      </View>
    );
  }
}

const ModalDatePicker: React.SFC<Props> = props => {
  return <StatefulModalDatePicker {...props} />;
}
export default ModalDatePicker;

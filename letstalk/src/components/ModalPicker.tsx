import React, { ReactNode } from 'react';
import {
  Dimensions,
  Modal,
  PickerIOS,
  PickerProperties,
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

type Props = PickerProperties & WrappedFieldProps & {
  label: string;
  children?: ReactNode;
};

interface State {
  modalVisible: boolean;
  values: Array<any>;
  labels: Array<string>;
};

class StatefulModalPicker extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const { values, labels } = this.getLabelsAndValues();

    this.state = {
      modalVisible: false,
      values,
      labels,
    };

    this.show = this.show.bind(this);
    this.hide = this.hide.bind(this);
    this.hideAndBlur = this.hideAndBlur.bind(this);
    this.onSubmitPress = this.onSubmitPress.bind(this);
  }

  getLabelsAndValues () {
    return {
      labels: React.Children.map(this.props.children, child => (child as any).props.label),
      values: React.Children.map(this.props.children, child => (child as any).props.value),
    };
  }

  componentWillReceiveProps(props: Props) {
    this.setState(this.getLabelsAndValues());
  }

  onSubmitPress() {
    // If the user presses submit without changing options, we choose first.
    const { value, onChange } = this.props.input;
    onChange(value || this.state.values[0]);
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
    const { label } = this.props;
    const { onChange, value } = this.props.input;
    const { error, touched, warning } = this.props.meta;
    const { modalVisible, values, labels } = this.state;
    const valueLabel = value ? labels[values.indexOf(value)] : label;

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
                <PickerIOS
                  style={styles.bottomPicker}
                  selectedValue={value}
                  onValueChange={onChange}
                >
                  {this.props.children}
                </PickerIOS>
              </View>
            </View>
          </View>
        </Modal>
        <ActionButton
          title={valueLabel}
          onPress={this.show}
        />
        {touched && (
          (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
          (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
      </View>
    );
  }
}

const ModalPicker: React.SFC<Props> = props => {
  return <StatefulModalPicker {...props} />;
}
export default ModalPicker;

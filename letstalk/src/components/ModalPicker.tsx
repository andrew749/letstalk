// NB: This is meant for iOS only.
// TODO: Rename to .ios.tsx
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
import BottomModal from './BottomModal';

const SCREEN_WIDTH = Dimensions.get('window').width;

const styles = StyleSheet.create({
  bottomPicker: {
    width: SCREEN_WIDTH,
  },
});

type Props = PickerProperties & WrappedFieldProps & {
  label: string;
  children?: ReactNode;
};

interface State {
  values: Array<any>;
  labels: Array<string>;
};

class StatefulModalPicker extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = this.getLabelsAndValues(props);
    this.setValueIfSingle();
  }

  getLabelsAndValues (props: Props) {
    return {
      labels: React.Children.map(props.children, child => (child as any).props.label),
      values: React.Children.map(props.children, child => (child as any).props.value),
    };
  }

  componentWillReceiveProps(props: Props) {
    this.setState(this.getLabelsAndValues(props), this.setValueIfSingle);
  }

  setValueIfSingle() {
    const { onChange, value } = this.props.input;
    if (this.state.values.length === 1 && this.state.values[0] !== value) {
      onChange(this.state.values[0]);
    }
  }

  render() {
    const { label } = this.props;
    const { onChange, value } = this.props.input;
    const { values, labels } = this.state;
    const valueLabel = value ? labels[values.indexOf(value)] : null;
    const onSubmitPress = () => {
      onChange(value || this.state.values[0]);
    }

    // TODO: Maybe hold state about what the value is using another onChange, and only call the
    // passed in onChange when the user presses submit.
    return (
      <BottomModal {...this.props} onSubmitPress={onSubmitPress} valueLabel={valueLabel}>
        <PickerIOS
          style={styles.bottomPicker}
          selectedValue={value}
          onValueChange={onChange}
        >
          {this.props.children}
        </PickerIOS>
      </BottomModal>
    );
  }
}

const ModalPicker: React.SFC<Props> = props => {
  return <StatefulModalPicker {...props} />;
}
export default ModalPicker;

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
import BottomModal from './BottomModal';

type Props = WrappedFieldProps & {
  label: string;
  defaultDate: Date;
  mode?: 'date' | 'time' | 'datetime';
};

const SCREEN_WIDTH = Dimensions.get('window').width;

const styles = StyleSheet.create({
  bottomPicker: {
    width: SCREEN_WIDTH,
  },
});

const ModalDatePicker: React.SFC<Props> = (props) => {
  const { defaultDate, label, mode } = props;
  const { onChange, value } = props.input;
  const onSubmitPress = () => {
    onChange(value || defaultDate);
  }
  // TODO: make this externally configurable
  const options = { year: 'numeric', month: 'long', day: 'numeric' };
  const valueLabel = value ? value.toLocaleDateString('en-US', options) : null;

  // TODO: Maybe hold state about what the value is using another onChange, and only call the
  // passed in onChange when the user presses submit.
  return (
    <BottomModal {...props} onSubmitPress={onSubmitPress} valueLabel={valueLabel}>
      <DatePickerIOS
        mode={mode}
        style={styles.bottomPicker}
        date={value || defaultDate}
        onDateChange={onChange}
      />
    </BottomModal>
  );
}
export default ModalDatePicker;

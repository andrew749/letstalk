// NB: This is meant for iOS only.
// TODO: Rename to .ios.tsx
import React, { ReactNode } from 'react';
import {
  Dimensions,
  Modal,
  Platform,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import DatePicker from 'react-native-datepicker';
import { WrappedFieldProps } from 'redux-form';
import { Button, ButtonProps, FormValidationMessage } from 'react-native-elements';
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
  // TODO: make this externally configurable
  const options = { year: 'numeric', month: 'long', day: 'numeric' };
  const valueLabel = value ? value.toLocaleDateString('en-US', options) : null;
  const pickerButtonStyle = {
    buttonStyle: {
      width: SCREEN_WIDTH - 40,
      borderRadius: 5,
    },
    textStyle: {
      color:'#000'
    }
  }

  // TODO: Maybe hold state about what the value is using another onChange, and only call the
  // passed in onChange when the user presses submit.

  const dateValue = (value || defaultDate);
  return (
    <DatePicker
        style={{width: '100%'}}
        date={ dateValue }
        mode="date"
        placeholder="Select Date"
        format="YYYY-MM-DD"
        minDate="1900-01-01"
        maxDate={new Date()}
        confirmBtnText="Confirm"
        cancelBtnText="Cancel"
        customStyles={{
          dateIcon: {
            position: 'absolute',
            left: 0,
            top: 4,
            marginLeft: 0
          },
          dateInput: {
            marginLeft: 36
          }
          // ... You can check the source to find the other keys.
        }}
        onDateChange={(date) => {
          onChange(new Date(date));
        }}
      />
  );
}
export default ModalDatePicker;

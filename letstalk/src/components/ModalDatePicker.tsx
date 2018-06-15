import React, { ReactNode } from 'react';
import {
  Dimensions,
  Modal,
  Platform,
  StyleProp,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableHighlight,
  TouchableWithoutFeedback,
  View,
  ViewStyle
} from 'react-native';
import DatePicker from 'react-native-datepicker';
import { WrappedFieldProps } from 'redux-form';
import { Button, ButtonProps, FormValidationMessage } from 'react-native-elements';
import BottomModal from './BottomModal';
import Moment from 'moment';
import Card from './Card';

type Props = WrappedFieldProps & {
  label: string;
  dateObj?: boolean; // whether to return a `Date` or a string
  defaultDate?: Date;
  mode?: 'date' | 'time' | 'datetime';
  cardStyle?: StyleProp<ViewStyle>
};

const SCREEN_WIDTH = Dimensions.get('window').width;

const styles = StyleSheet.create({
  datePickerButton: {
    borderWidth: 0,
  },
  label: {
    fontWeight: 'bold',
    fontSize: 14,
  },
  card: {
    flex: 1,
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 10,
  },
});

const ModalDatePicker: React.SFC<Props> = (props) => {
  const { defaultDate, label, mode, dateObj } = props;
  const { onChange, onBlur, value } = props.input;
  const { error, touched, warning } = props.meta;
  // TODO: make this externally configurable
  const options = { year: 'numeric', month: 'long', day: 'numeric' };
  const pickerButtonStyle = {
    buttonStyle: {
      width: SCREEN_WIDTH - 40,
      borderRadius: 5,
    },
    textStyle: {
      color:'#000',
    },
  }

  // TODO: Maybe hold state about what the value is using another onChange, and only call the
  // passed in onChange when the user presses submit.

  const dateValue = (value || defaultDate);
  return (
    <Card style={[styles.card, props.cardStyle]}>
      <Text style={styles.label}>{label}</Text>
      <DatePicker
        date={ dateValue }
        mode="date"
        showIcon={false}
        placeholder="Select Date"
        format="YYYY-MM-DD"
        minDate={new Date(1900, 1, 1)}
        maxDate={new Date()}
        confirmBtnText="Confirm"
        cancelBtnText="Cancel"
        // @ts-ignore waiting for https://github.com/DefinitelyTyped/DefinitelyTyped/pull/26237 to land
        getDateStr={(date: Date) => Moment(date).format("MMM Do, YYYY")}
        customStyles={{
          dateInput: styles.datePickerButton,
        }}
        onCloseModal={onBlur as () => void}
        onDateChange={(dateString, date) => {
          if (!dateObj) {
            onChange(Moment(date).format("YYYY-MM-DD"));
          } else {
            onChange(date);
          }
        }}
      />
      {touched && (
        (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
        (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
    </Card>
  );
}
export default ModalDatePicker;

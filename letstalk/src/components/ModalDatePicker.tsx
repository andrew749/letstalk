// NB: This is meant for iOS only.
// TODO: Rename to .ios.tsx
import React, { ReactNode } from 'react';
import {
  Dimensions,
  Modal,
  DatePickerIOS,
  DatePickerIOSProperties,
  DatePickerAndroid,
  Platform,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
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

const launchAndroidDatePicker = async (
    onSubmit: (newValue: Date) => void
): Promise<void> =>  {
  try {
    const {action, year, month, day} = await DatePickerAndroid.open({
      date: new Date()
    });
    if (action !== DatePickerAndroid.dismissedAction) {
      onSubmit(new Date(year, month, day));
    }
  } catch ({code, message}) {
    console.warn('Cannot open date picker', message);
  }
}

const ModalDatePicker: React.SFC<Props> = (props) => {
  const { defaultDate, label, mode } = props;
  const { onChange, value } = props.input;
  const onSubmitPress = () => {
    onChange(value || defaultDate);
  }
  // TODO: make this externally configurable
  const options = { year: 'numeric', month: 'long', day: 'numeric' };
  const valueLabel = value ? value.toLocaleDateString('en-US', options) : null;

  const pickerButtonStyle = {
    buttonStyle: {
      width: SCREEN_WIDTH - 40,
      borderRadius: 30,
    },
    textStyle: {
      color:'#000'
    }
  }

  // TODO: Maybe hold state about what the value is using another onChange, and only call the
  // passed in onChange when the user presses submit.

  const dateValue = (value || defaultDate).toString();
  const datePicker = Platform.select({
    'ios': (
    <BottomModal {...props}
      onSubmitPress={onSubmitPress}
      valueLabel={valueLabel}>
      <DatePickerIOS
        mode={mode}
        style={styles.bottomPicker}
        date={value || defaultDate}
        onDateChange={onChange}
        />
    </BottomModal>
    ),
    'android': (
      <Button
        {...props}
        title={dateValue}
        textStyle={pickerButtonStyle.textStyle}
        buttonStyle={pickerButtonStyle.buttonStyle}
        onPress={launchAndroidDatePicker.bind(this, onChange)}
      />
     )}
   );
  return datePicker;
}
export default ModalDatePicker;

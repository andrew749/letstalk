import React, { ReactNode } from 'react';
import {
  PickerProperties,
  Picker,
  View,
} from 'react-native';
import { FormValidationMessage } from 'react-native-elements';
import { WrappedFieldProps } from 'redux-form';

type Props = PickerProperties & WrappedFieldProps & {
  label: string;
  children?: ReactNode;
};

const AndroidModalPicker: React.SFC<Props> = props => {
  const { label } = props;
  const { onChange, value } = props.input;
  const { error, touched, warning } = props.meta;
  // TODO: Actually test this
  return (
    <View>
      <Picker {...props} onValueChange={onChange} selectedValue={value} prompt={label} />
      {touched && (
        (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
        (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
    </View>
  );
};

export default AndroidModalPicker;

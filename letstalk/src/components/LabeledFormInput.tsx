import React, { SFC } from 'react';
import { View } from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { FormValidationMessage, FormInputProps, FormInput, FormLabel } from 'react-native-elements';

type Props = FormInputProps & WrappedFieldProps & { label: string }

const LabeledFormInput: SFC<Props> = props => {
  const { label } = props;
  const { onChange, onBlur, value } = props.input;
  const { error, touched, warning } = props.meta;
  return (
    <View>
      {label && <FormLabel>{label}</FormLabel>}
      <FormInput
        {...props}
        onBlur={onBlur as () => void} // Thanks jhang (type hack to make this typecheck)
        onChangeText={onChange}
        value={value}
      />
      {touched && (
        (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
        (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
    </View>
  );
};

export default LabeledFormInput;


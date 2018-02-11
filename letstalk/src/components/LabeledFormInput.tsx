import React, { Component } from 'react';
import { View } from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { FormValidationMessage, FormInputProps, FormInput, FormLabel } from 'react-native-elements';

type Props = FormInputProps & WrappedFieldProps & { label: string }

const LabeledFormInput: React.SFC<Props> = props => {
  const formLabel = props.label === null ? null : <FormLabel>{props.label}</FormLabel>;
  const { onChange, onBlur, value } = props.input;
  const { error, touched, warning } = props.meta;
  console.log(touched);
  return (
    <View>
      {formLabel}
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


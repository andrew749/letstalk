import React, { Component } from 'react';
import { View } from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { FormInputProps, FormInput, FormLabel } from 'react-native-elements';

type Props = FormInputProps & WrappedFieldProps & { label: string }

const LabeledFormInput: React.SFC<Props> = props => {
  const formLabel = props.label === null ? null : <FormLabel>{props.label}</FormLabel>;
  const { onChange, value } = props.input;
  return (
    <View>
      {formLabel}
      <FormInput
        {...props}
        onChangeText={onChange}
        value={value}
      />
    </View>
  );
};

export default LabeledFormInput;


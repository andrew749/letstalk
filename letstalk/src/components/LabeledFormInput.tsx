import React, { SFC } from 'react';
import { View } from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { FormValidationMessage, FormInputProps, FormInput, FormLabel } from 'react-native-elements';

type Props = FormInputProps & WrappedFieldProps & { label: string, onSubmitEditing?: () => void }

class LabeledFormInput extends React.Component<Props> {
  private elementRef: React.Ref<FormInput>;
  constructor(props: Props) {
    super(props);
    // @ts-ignore
    this.elementRef = React.createRef();
  }
  focus() {
    // @ts-ignore
    this.elementRef.current.focus();
  }
  render() {
    const props = this.props;
    const { label } = props;
    const { onChange, onBlur, value } = props.input;
    const { error, touched, warning } = props.meta;
    return (
      <View>
        {label && <FormLabel>{label}</FormLabel>}
        <FormInput
          {...props}
          ref={this.elementRef}
          onBlur={onBlur as () => void} // Thanks jhang (type hack to make this typecheck)
          onChangeText={onChange}
          value={value}
          onSubmitEditing={props.onSubmitEditing}
        />
        {touched && (
          (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
          (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
        </View>
      );
    }
};

export default LabeledFormInput;

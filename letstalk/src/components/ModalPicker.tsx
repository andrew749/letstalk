import React, { ReactNode } from 'react';
import {
  PickerProperties,
  View,
} from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';
import ActionButton from './ActionButton';
import SimplePicker from 'react-native-simple-picker';

type Props = PickerProperties & WrappedFieldProps & {
  label: string;
  children?: ReactNode;
};

interface State {
  values: Array<any>;
  labels: Array<string>;
};

const getLabelsAndValues = (children: ReactNode) => {
  return {
    labels: React.Children.map(children, child => (child as any).props.label),
    values: React.Children.map(children, child => (child as any).props.value),
  };
}

class StatefulModalPicker extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = getLabelsAndValues(props.children);
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setState(getLabelsAndValues(nextProps.children));
  }

  render() {
    const { label } = this.props;
    const { onBlur, onChange, value } = this.props.input;
    const { error, touched, warning } = this.props.meta;
    const { values, labels } = this.state;
    const valueLabel = value ? labels[values.indexOf(value)] : label;
    // TODO: Don't use ActionButton
    // Also, might want to just copy the code from SimplePicker here, since we need the onClose
    // action to forward to onBlur. Also, will make it easier to style.
    return (
      <View>
        <ActionButton
          title={valueLabel || label}
          onPress={() => (this.refs.picker as any).show()}
        />
        <SimplePicker
          ref={'picker'}
          options={values}
          labels={labels}
          onSubmit={(value: any) => {onChange(value); (onBlur as () => void)()}}
        />
        {touched && (
          (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
          (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
      </View>
    );
  }
}

const ModalPicker: React.SFC<Props> = props => {
  return <StatefulModalPicker {...props} />;
}
export default ModalPicker;

import React, { ReactNode } from 'react';
import {
  Dimensions,
  Modal,
  PickerIOS,
  PickerProperties,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import { ButtonGroup, FormValidationMessage } from 'react-native-elements';
import Card from './Card';

const styles = StyleSheet.create({
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

type Props = PickerProperties & WrappedFieldProps & {
  label: string;
  children?: ReactNode;
};

interface State {
  values: Array<any>;
  labels: Array<string>;
};

class StatefulButtonPicker extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const { values, labels } = this.getLabelsAndValues();

    this.state = {
      values,
      labels,
    };
  }

  getLabelsAndValues () {
    return {
      labels: React.Children.map(this.props.children, child => (child as any).props.label),
      values: React.Children.map(this.props.children, child => (child as any).props.value),
    };
  }

  componentWillReceiveProps(props: Props) {
    this.setState(this.getLabelsAndValues());
  }

  render() {
    const { label } = this.props;
    const { onChange, value } = this.props.input;
    const { error, touched, warning } = this.props.meta;
    const { values, labels } = this.state;
    const valueLabel = value ? labels[values.indexOf(value)] : null;
    const selectedIndex = value ? values.indexOf(value) : -1;

    return (
      <Card style={styles.card}>
        <Text style={styles.label}>{label}</Text>
        <ButtonGroup
          onPress={(i: number) => onChange(values[i])}
          selectedIndex={selectedIndex}
          buttons={labels}
        />
        {touched && (
          (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
          (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
      </Card>
    );
  }
}

const ButtonPicker: React.SFC<Props> = props => {
  return <StatefulButtonPicker {...props} />;
}
export default ButtonPicker;

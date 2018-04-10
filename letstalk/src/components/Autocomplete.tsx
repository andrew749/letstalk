import Immutable from 'immutable';
import fuzzysearch from 'fuzzysearch';

import React, { Component } from 'react';
import {
  ListView,
  ListViewDataSource,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';

interface FilterableElement {
  readonly id: any;
  readonly value: string;
}

interface Props {
  data: Immutable.List<FilterableElement>;
}

interface State {
  curValue: string;
  filteredElements: Immutable.List<FilterableElement>;
}

function filterElements(
  value: string,
  elements: Immutable.List<FilterableElement>,
): Immutable.List<FilterableElement> {
  const lowerVal = value.toLowerCase();
  return elements.filter((elem: FilterableElement) => {
    return fuzzysearch(lowerVal, elem.value.toLowerCase());
  }).toList();
}

class Autocomplete extends Component<Props, State> {
  private ds: ListViewDataSource;

  constructor(props: Props) {
    super(props);

    this.ds = new ListView.DataSource({
      rowHasChanged: (r1: FilterableElement, r2: FilterableElement) => r1.id !== r2.id
    });

    this.state = {
      curValue: '',
      filteredElements: props.data,
    };

    this.filterElements = this.filterElements.bind(this);
  }

  componentWillReceiveProps(props: Props) {
    this.setState({
      filteredElements: filterElements(this.state.curValue, props.data),
    });
  }

  filterElements(newValue: string) {
    this.setState({
      curValue: newValue,
      filteredElements: filterElements(newValue, this.props.data),
    });
  }

  renderElement(elem: FilterableElement) {
    return (
      <Text>{elem.value}</Text>
    );
  }

  render() {
    const { filteredElements } = this.state;
    const ds = this.ds.cloneWithRows(filteredElements.toJS());
    return (
      <View style={styles.container}>
        <TextInput
            style={styles.textinput}
            onChangeText={this.filterElements}
            placeholder="Looking for a friend with..."
        />
        <ListView
          dataSource={ds}
          renderRow={this.renderElement}
        />
      </View>
    );
  }
}

export default Autocomplete;

var styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  textinput: {
    marginTop: 30,
    height: 40,
    width: 200
  }
});

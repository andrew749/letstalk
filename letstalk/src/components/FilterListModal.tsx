import Immutable from 'immutable';
import fuzzysearch from 'fuzzysearch';

import React, { Component } from 'react';
import {
  ListView,
  ListViewDataSource,
  Modal,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

export interface FilterableElement {
  readonly id: any;
  readonly value: string;
}

interface Props {
  data: Immutable.List<FilterableElement>;
  onSelect(elem: FilterableElement): Promise<void>;
  placeholder: string;
}

interface State {
  curValue: string;
  filteredElements: Immutable.List<FilterableElement>;
  modalVisible: boolean;
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

class FilterListModal extends Component<Props, State> {
  private ds: ListViewDataSource;
  private textInputRef: TextInput;

  constructor(props: Props) {
    super(props);

    this.ds = new ListView.DataSource({
      rowHasChanged: (r1: FilterableElement, r2: FilterableElement) => r1.id !== r2.id
    });

    this.state = {
      curValue: '',
      filteredElements: props.data,
      modalVisible: false,
    };

    this.filterElements = this.filterElements.bind(this);
    this.renderElement = this.renderElement.bind(this);
    this.setModalVisible = this.setModalVisible.bind(this);
    this.textInputRef = null;
  }

  componentWillReceiveProps(props: Props) {
    this.setState({
      filteredElements: filterElements(this.state.curValue, props.data),
    });
  }

  private setModalVisible(modalVisible: boolean) {
    this.setState({ modalVisible });
  }

  private filterElements(newValue: string) {
    this.setState({
      curValue: newValue,
      filteredElements: filterElements(newValue, this.props.data),
    });
  }

  private renderElement(elem: FilterableElement) {
    const onPress = async () => {
      await this.props.onSelect(elem);
      this.setState({
        curValue: '',
        filteredElements: this.props.data,
        modalVisible: false,
      });
    };
    return (
      <TouchableOpacity style={styles.item} onPress={onPress}>
        <Text style={styles.itemText}>{elem.value}</Text>
      </TouchableOpacity>
    );
  }

  render() {
    const { placeholder } = this.props;
    const { filteredElements } = this.state;
    const ds = this.ds.cloneWithRows(filteredElements.toJS());
    return (
      <View>
        <View style={styles.textInputContainer}>
          <TouchableOpacity
            onPress={() => { this.setModalVisible(true) }}
          >
            <Text style={styles.buttonText}>{placeholder}</Text>
          </TouchableOpacity>
        </View>
        <Modal
          animationType="fade"
          transparent={false}
          visible={this.state.modalVisible}
          onShow={() => this.textInputRef.focus()}
        >
          <View style={styles.container}>
            <View style={styles.textInputContainer}>
              <TextInput
                style={styles.textInput}
                onChangeText={this.filterElements}
                onFocus={() => this.setModalVisible(true)}
                placeholder={placeholder}
                ref={textInput => this.textInputRef = textInput as any}
              />
            </View>
            <ListView
              keyboardShouldPersistTaps={'always'}
              dataSource={ds}
              renderRow={this.renderElement}
            />
          </View>
        </Modal>
      </View>
    );
  }
}

export default FilterListModal;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 25,
    justifyContent: 'center',
    backgroundColor: '#F0F0F0',
  },
  item: {
    flex: 1,
    justifyContent: 'center',
    borderTopWidth: 0.5,
    borderTopColor: '#909090',
    padding: 10,
  },
  itemText: {
    fontSize: 18,
    fontWeight: '700',
  },
  textInputContainer: {
    justifyContent: 'center',
    margin: 10,
    borderRadius: 25,
    backgroundColor: 'white',
    padding: 10,
    paddingLeft: 20,
    paddingRight: 20,
    height: 50,
  },
  buttonText: {
    fontSize: 18,
    color: '#909090',
  },
  textInput: { },
});

import Immutable from 'immutable';
import fuzzysearch from 'fuzzysearch';
import { MaterialIcons } from '@expo/vector-icons';

import React, { Component } from 'react';
import {
  Dimensions,
  ListView,
  ListViewDataSource,
  Modal,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

const SCREEN_WIDTH = Dimensions.get('window').width;

export interface FilterableElement {
  readonly id: any;
  readonly value: string;
}

interface FilterableElementType extends FilterableElement {
  readonly type: 'FILTERABLE_ELEMENT',
}

interface GapType extends FilterableElement {
  readonly type: 'GAP',
}

interface NoMoreResultsType extends FilterableElement {
  readonly type: 'NO_MORE_RESULTS',
}

type Element = FilterableElementType | GapType | NoMoreResultsType;

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

  private renderElement(elem: Element) {
    console.log(elem);
    switch (elem.type) {
      case 'FILTERABLE_ELEMENT':
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
      case 'GAP':
        return <View style={styles.gap} />;
      case 'NO_MORE_RESULTS':
        return (
          <View style={styles.noMoreResults}>
            <Text>No more results...</Text>
          </View>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = elem;
    }
  }

  render() {
    const { placeholder } = this.props;
    const { filteredElements } = this.state;
    const elements = filteredElements.map(elem => {
      return {...elem, type: 'FILTERABLE_ELEMENT'};
    }).toJS();
    const allItems = [{ type: 'GAP' }].concat(elements).concat([{ type: 'NO_MORE_RESULTS' }]);
    const ds = this.ds.cloneWithRows(allItems);
    const onPressDismiss = () => {
      this.setState({
        curValue: '',
        filteredElements: this.props.data,
        modalVisible: false,
      });
    };
    return (
      <View>
        <View style={styles.textInputContainer}>
          <TouchableOpacity onPress={() => { this.setModalVisible(true) }}>
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
            <View style={styles.topContainer}>
              <View style={[styles.textInputContainer, styles.textInputContainerWidth]}>
                <TextInput
                  style={styles.textInput}
                  onChangeText={this.filterElements}
                  onFocus={() => this.setModalVisible(true)}
                  placeholder={placeholder}
                  placeholderTextColor='white'
                  ref={textInput => this.textInputRef = textInput as any}
                />
              </View>
              <TouchableOpacity style={styles.dismiss} onPress={onPressDismiss}>
                <MaterialIcons name="close" size={24} />
              </TouchableOpacity>
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

const DISMISS_BUTTON_PADDING = 12;
// Size of button + 2 * button padding + my margin
const TEXT_INPUT_RIGHT_MARGIN = 24 + 2 * DISMISS_BUTTON_PADDING + 10;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 25,
    justifyContent: 'center',
    backgroundColor: 'white',
  },
  item: {
    flex: 1,
    justifyContent: 'center',
    borderBottomWidth: 0.5,
    borderColor: '#909090',
    backgroundColor: '#F0F0F0',
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
    backgroundColor: '#FFD475',
    paddingLeft: 10,
    paddingRight: 10,
    height: 30,
  },
  textInputContainerWidth: {
    marginRight: 0,
    width: SCREEN_WIDTH - TEXT_INPUT_RIGHT_MARGIN,
  },
  dismiss: {
    flex: 0,
    padding: 12,
  },
  topContainer: {
    justifyContent: 'space-between',
    flexDirection: 'row',
    width: SCREEN_WIDTH,
    backgroundColor: '#FFBB25',
  },
  buttonText: {
    fontSize: 14,
    color: 'white',
  },
  gap: {
    height: 20,
    borderBottomWidth: 0.5,
    borderColor: '#909090',
  },
  textInput: {
    fontSize: 14,
    color: 'white',
  },
  noMoreResults: {
    marginTop: 10,
    fontSize: 14,
    alignItems: 'center',
  },
});

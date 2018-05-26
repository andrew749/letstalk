import Immutable from 'immutable';
import fuzzysearch from 'fuzzysearch';
import Fuse from 'fuse-js-latest';
import { MaterialIcons } from '@expo/vector-icons';

import React, { Component, ReactNode } from 'react';
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

import {
  SearchBar,
} from 'react-native-elements';

const SCREEN_WIDTH = Dimensions.get('window').width;

export interface FilterableElement {
  readonly id: any;
  readonly value: string;
}

interface FilterableElementType extends FilterableElement {
  readonly type: 'FILTERABLE_ELEMENT';
  readonly searchValue: string;
}

interface GapType {
  readonly type: 'GAP';
}

interface NoMoreResultsType {
  readonly type: 'NO_MORE_RESULTS';
}

interface RawInputType {
  readonly type: 'RAW_INPUT';
  readonly searchValue: string;
}

type Element = FilterableElementType | GapType | NoMoreResultsType | RawInputType;

interface Props {
  data: Immutable.List<FilterableElement>;
  onSelect(elem: FilterableElement): Promise<void>;
  placeholder: string;
  onRawSelect?(value: string): Promise<void>;
  buttonComponent?(onPress: () => void): ReactNode;
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
  if (value === '') return elements;

  const options = {
    shouldSort: true,
    threshold: 0.6,
    location: 0,
    distance: 100,
    maxPatternLength: 32,
    minMatchCharLength: 1,
  };
  const fuse = new Fuse(elements.map(elem => elem.value.toLowerCase()).toJS(), options);
  const result = fuse.search(value.toLowerCase());

  return Immutable.List(result.map((idx: number) => elements.get(idx)));
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
    let onPress;
    switch (elem.type) {
      case 'FILTERABLE_ELEMENT':
        onPress = async () => {
          await this.props.onSelect(elem);
          this.setState({
            curValue: '',
            filteredElements: this.props.data,
            modalVisible: false,
          });
        };
        const options = {
          includeMatches: true,
          threshold: 0.6,
          location: 0,
          distance: 100,
          maxPatternLength: 32,
          minMatchCharLength: 1,
        };

        const lowerVal = elem.searchValue;
        const fuse = new Fuse([elem.value.toLowerCase()], options);
        const result: any = fuse.search(lowerVal);
        if (result.length > 0 && result[0].matches.length > 0 &&
          result[0].matches[0].indices.length > 0
        ) {
          const [start_idx, end_idx] = result[0].matches[0].indices.reduce(
            (acc: Array<number>, cur: Array<number>) => {
            if (cur[1] - cur[0] > acc[1] - acc[0]) return cur;
            else return acc;
          });


          if (start_idx < end_idx && end_idx < elem.value.length) {
            const start = elem.value.slice(0, start_idx);
            const middle = elem.value.slice(start_idx, end_idx+1);
            const end = elem.value.slice(end_idx+1);
            console.log({ start_idx, end_idx, start, middle, end });
            return (
              <TouchableOpacity style={styles.item} onPress={onPress}>
                <Text style={styles.itemText}>
                  {start}
                  <Text style={{fontWeight: 'bold'}}>{middle}</Text>
                  {end}
                </Text>
              </TouchableOpacity>
            );
          }
        }

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
      case 'RAW_INPUT':
        onPress = async () => {
          await this.props.onRawSelect(elem.searchValue);
          this.setState({
            curValue: '',
            filteredElements: this.props.data,
            modalVisible: false,
          });
        };
        return (
          <TouchableOpacity style={styles.item} onPress={onPress}>
            <Text style={[styles.itemText, styles.rawInputText]}>
              {'Add '}
              <Text style={{fontWeight: 'bold'}}>"{elem.searchValue}"</Text>
            </Text>
          </TouchableOpacity>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = elem;
    }
  }

  render() {
    const { placeholder, onRawSelect, buttonComponent } = this.props;
    const { filteredElements, curValue } = this.state;
    let elements = filteredElements.map(elem => {
      return {...elem, type: 'FILTERABLE_ELEMENT', searchValue: curValue };
    }).toJS();
    if (filteredElements.size > 0) elements = [{ type: 'GAP' }].concat(elements);
    if (!!onRawSelect && curValue !== '') {
      elements = elements.concat([{ type: 'GAP' }, { type: 'RAW_INPUT', searchValue: curValue }]);
    }
    elements = elements.concat([{ type: 'NO_MORE_RESULTS' }]);
    const ds = this.ds.cloneWithRows(elements);
    const onPressDismiss = () => {
      this.setState({
        curValue: '',
        filteredElements: this.props.data,
        modalVisible: false,
      });
    };

    const onPressOpen = () => {
      this.setModalVisible(true)
    }

    const butt = !!buttonComponent ? buttonComponent(onPressOpen) : (
      <TouchableOpacity
        style={styles.textInputContainer}
        onPress={onPressOpen}
      >
        <MaterialIcons name="search" size={20} color="white" />
        <Text style={styles.buttonText}>{placeholder}</Text>
      </TouchableOpacity>
    );

    return (
      <View>
        { butt }
        <Modal
          animationType="fade"
          transparent={true}
          visible={this.state.modalVisible}
          onShow={() => this.textInputRef.focus()}
        >
          <View style={styles.container}>
            <View style={styles.topContainer}>
              <View style={[styles.textInputContainer, styles.textInputContainerWidth]}>
                <MaterialIcons name="search" size={20} color="white" />
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
    marginTop: 20,
    justifyContent: 'center',
    backgroundColor: 'white',
  },
  item: {
    justifyContent: 'center',
    borderBottomWidth: 0.5,
    borderColor: '#909090',
    padding: 10,
  },
  rawInputText: {
    color: '#003CB2',
  },
  itemText: {
    fontSize: 18,
  },
  textInputContainer: {
    alignItems: 'center',
    flexDirection: 'row',
    margin: 10,
    borderRadius: 25,
    backgroundColor: '#FFD475',
    paddingLeft: 10,
    paddingRight: 20,
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
    backgroundColor: '#FFC107',
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
    alignItems: 'center',
  },
});

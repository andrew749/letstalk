import Immutable from 'immutable';
import fuzzysearch from 'fuzzysearch';
import Fuse from 'fuse-js-latest';
import { MaterialIcons } from '@expo/vector-icons';

import React, { Component, ReactNode } from 'react';
import {
  Dimensions,
  EmitterSubscription,
  Keyboard,
  ListView,
  ListViewDataSource,
  Modal,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { Constants } from 'expo';

import {
  SearchBar,
} from 'react-native-elements';

const SCREEN_WIDTH = Dimensions.get('window').width;
const SCREEN_HEIGHT = Dimensions.get('window').height;

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
  onRawSelect?(value: string): Promise<void>;
  curValue: string;
}

interface State {
  filteredElements: Immutable.List<FilterableElement>;
  keyboardHeight: number;
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
  private keyboardDidShowListener: EmitterSubscription;
  private keyboardDidHideListener: EmitterSubscription;

  constructor(props: Props) {
    super(props);

    this.ds = new ListView.DataSource({
      rowHasChanged: (r1: FilterableElement, r2: FilterableElement) => r1.id !== r2.id
    });

    this.state = {
      filteredElements: props.data,
      keyboardHeight: 0,
    };

    this.renderElement = this.renderElement.bind(this);
    this.keyboardDidHide = this.keyboardDidHide.bind(this);
    this.keyboardDidShow = this.keyboardDidShow.bind(this);
  }

  componentWillReceiveProps(props: Props) {
    this.setState({
      filteredElements: filterElements(props.curValue, props.data),
    });
  }

  keyboardDidShow (e: any) {
    let newSize = e.endCoordinates.height;
    console.log(newSize);
    this.setState({ keyboardHeight: newSize });
  }

  keyboardDidHide (e: any) {
    this.setState({ keyboardHeight: 0 });
  }

  componentWillMount () {
    this.keyboardDidShowListener = Keyboard.addListener(
      'keyboardDidShow', this.keyboardDidShow,
    );
    this.keyboardDidHideListener = Keyboard.addListener(
      'keyboardDidHide', this.keyboardDidHide,
    );
  }

  componentWillUnmount () {
    this.keyboardDidShowListener.remove()
    this.keyboardDidHideListener.remove()
  }

  private renderElement(elem: Element) {
    let onPress;
    switch (elem.type) {
      case 'FILTERABLE_ELEMENT':
        onPress = async () => {
          await this.props.onSelect(elem);
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
        const padding = {
          paddingBottom: this.state.keyboardHeight + 80,
        };
        return (
          <View style={[styles.noMoreResults, padding]}>
            <Text>No more results...</Text>
          </View>
        );
      case 'RAW_INPUT':
        onPress = async () => {
          await this.props.onRawSelect(elem.searchValue);
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
    const { curValue, onRawSelect } = this.props;
    const { filteredElements } = this.state;
    let elements = filteredElements.map(elem => {
      return {...elem, type: 'FILTERABLE_ELEMENT', searchValue: curValue };
    }).toJS();
    if (filteredElements.size > 0) elements = [{ type: 'GAP' }].concat(elements);
    if (!!onRawSelect && curValue !== '') {
      elements = elements.concat([{ type: 'GAP' }, { type: 'RAW_INPUT', searchValue: curValue }]);
    }
    elements = elements.concat([{ type: 'NO_MORE_RESULTS' }]);
    const ds = this.ds.cloneWithRows(elements);

    return (
      <View style={styles.container}>
        <ListView
          keyboardShouldPersistTaps={'always'}
          dataSource={ds}
          renderRow={this.renderElement}
        />
      </View>
    );
  }
}

export default FilterListModal;

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    width: SCREEN_WIDTH,
    // TODO: doesn't take into account header height
    height: SCREEN_HEIGHT - Constants.statusBarHeight,
    top: 0,
    left: 0,
    alignSelf: 'stretch',
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
  gap: {
    height: 20,
    borderBottomWidth: 0.5,
    borderColor: '#909090',
  },
  noMoreResults: {
    marginTop: 10,
    alignItems: 'center',
  },
});

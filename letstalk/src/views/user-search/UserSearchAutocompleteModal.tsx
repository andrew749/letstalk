import React, { Component, ReactNode } from 'react';
import Immutable from 'immutable';
import fuzzysearch from 'fuzzysearch';
import Fuse from 'fuse-js-latest';
import { MaterialIcons } from '@expo/vector-icons';
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

import { MultiTrait } from '../../models/multi-trait';

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

interface HintType {
  readonly type: 'HINT';
  readonly hintText: string;
}

type Element = FilterableElementType | GapType | NoMoreResultsType | RawInputType | HintType;

interface Props {
  data: Immutable.List<MultiTrait>;
  onSelect(trait: MultiTrait): Promise<void>;
  value: string;
}

interface State {
  keyboardHeight: number;
}

class UserSearchAutocompleteModal extends Component<Props, State> {
  private ds: ListViewDataSource;
  private keyboardDidShowListener: EmitterSubscription;
  private keyboardDidHideListener: EmitterSubscription;

  constructor(props: Props) {
    super(props);

    this.ds = new ListView.DataSource({
      rowHasChanged: (r1: FilterableElement, r2: FilterableElement) => r1.id !== r2.id
    });

    this.state = {
      keyboardHeight: 0,
    };

    this.renderElement = this.renderElement.bind(this);
    this.keyboardDidHide = this.keyboardDidHide.bind(this);
    this.keyboardDidShow = this.keyboardDidShow.bind(this);
  }

  keyboardDidShow (e: any) {
    let newSize = e.endCoordinates.height;
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

  private renderElement(trait: MultiTrait) {
    // TODO: Render different kinds of elements.
    return <Text>yo</Text>;
  }

  render() {
    const { data, value } = this.props;
    const ds = this.ds.cloneWithRows(data.toJS());

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

export default UserSearchAutocompleteModal;

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
  hint: {
    paddingTop: 15,
    paddingLeft: 5,
    paddingRight: 5,
    fontSize: 14,
    alignItems: 'center',
  },
});

import React, { Component, ReactNode } from 'react';
import Immutable from 'immutable';
import { MaterialIcons } from '@expo/vector-icons';
import {
  EmitterSubscription,
  Keyboard,
  ListView,
  ListViewDataSource,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import Constants from 'expo-constants';

import {
  MultiTrait,
  MultiTraitTypes,
} from '../../models/multi-trait';
import Colors from '../../services/colors';
import Window from '../../services/window';

interface ItemElement {
  readonly type: 'ITEM';
  readonly id: string;
  readonly trait: MultiTrait;
}

interface EndOfItemsElement {
  readonly type: 'END_OF_ITEMS';
  readonly id: 'END_OF_ITEMS';
  readonly numElems: number;
}

interface HintElement {
  readonly type: 'HINT';
  readonly id: 'HINT';
  readonly hintText: string;
}

type Element = ItemElement | EndOfItemsElement | HintElement;

const CLICK_HINT: HintElement = {
  id: 'HINT',
  type: 'HINT',
  hintText: 'Click on an option to search for users',
};

interface Props {
  data: Immutable.List<MultiTrait>;
  onSelect(trait: MultiTrait): Promise<void>;
  value: string;
  errorMsg: string;
}

interface State {
  keyboardHeight: number;
}

function longestMatch(str: string, query: string) {
  str = str.toLowerCase();
  query = query.toLowerCase();
  let longestMatchIdxs = null;
  for (let i = 0; i < str.length; i++) {
    for (let j = 0; i + j < str.length && j < query.length; j++) {
      if (query[j] === str[i+j]) {
        if (longestMatchIdxs === null || j + 1 > (longestMatchIdxs[1] - longestMatchIdxs[0])) {
          longestMatchIdxs = [i, i+j+1];
        }
      } else {
        break;
      }
    }
  }
  return longestMatchIdxs;
}

class UserSearchAutocompleteModal extends Component<Props, State> {
  private ds: ListViewDataSource;
  private keyboardDidShowListener: EmitterSubscription;
  private keyboardDidHideListener: EmitterSubscription;

  constructor(props: Props) {
    super(props);

    this.ds = new ListView.DataSource({
      rowHasChanged: (r1: Element, r2: Element) => r1.id !== r2.id
    });

    this.state = {
      keyboardHeight: 0,
    };

    this.renderElement = this.renderElement.bind(this);
    this.renderItem = this.renderItem.bind(this);
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

  private renderItem(trait: MultiTrait) {
    let icon: string = 'color-lens';
    let color: string = Colors.HIVE_PRIMARY;
    let tpe: string = 'Trait';
    let text = null;

    const name = trait.traitName;

    const matched = longestMatch(name, this.props.value);
    if (!!matched) {
      text = (
        <Text>
          {name.substring(0, matched[0])}
          <Text style={styles.boldText}>{name.substring(matched[0], matched[1])}</Text>
          {name.substring(matched[1])}
        </Text>
      );
    } else {
      text = <Text>{name}</Text>;
    }

    switch (trait.traitType) {
      case MultiTraitTypes.COHORT:
        icon = 'school';
        color = Colors.HIVE_SUBDUED;
        tpe = 'Cohort';
        break;
      case MultiTraitTypes.POSITION:
        icon = 'supervisor-account';
        color = Colors.HIVE_ACCENT;
        tpe = 'Position';
        break;
      case MultiTraitTypes.GROUP:
        icon = 'group-work';
        color = Colors.GROUP_GREEN;
        tpe = 'Group';
        break;
    }

    const onPress = () => this.props.onSelect(trait);

    return (
      <TouchableOpacity style={styles.searchItem} onPress={onPress}>
        <MaterialIcons name={icon} color={color} size={24} />
        <View style={styles.traitContainer}>
          <Text style={styles.traitText}>{text}</Text>
          <Text style={[styles.typeText, { color }]}>{tpe}</Text>
        </View>
      </TouchableOpacity>
    );
  }

  private renderElement(el: Element) {
    const { value, data, errorMsg } = this.props;

    switch (el.type) {
      case 'ITEM':
        return this.renderItem(el.trait);
      case 'END_OF_ITEMS':
        const padding = {
          paddingBottom: this.state.keyboardHeight + 80,
        };
        let text = '';
        if (value === '' && data.isEmpty()) {
          text = 'Start typing to see search options';
        } else if (data.isEmpty()) {
          text = 'No search options found for ' + value + ' - try something else';
        } else if (data.size === 1) {
          text = 'Showing the best search option';
        } else {
          text = 'Showing the best ' + el.numElems + ' search options';
        }

        const errorElem = !!errorMsg ? (
          <Text style={styles.errorText}>{ errorMsg }</Text>
        ) : null;

        return (
          <View style={[styles.noMoreResults, padding]}>
            <Text style={{color: 'gray'}}>{ text }</Text>
            {errorElem}
          </View>
        );
      case 'HINT':
        return (
          <View style={[styles.hint]}>
            <Text style={{color: 'gray'}}>{ el.hintText }</Text>
          </View>
        )
      default:
        // Ensure exhaustiveness of select
        const _: never = el;
    }
  }

  render() {
    const { data, value } = this.props;
    let els: Immutable.List<Element> = data.isEmpty() ?
      Immutable.List() : Immutable.List([CLICK_HINT])
    els = els.concat(data.map(trait => {
      const el: ItemElement = {
        type: 'ITEM',
        id: trait.traitName,
        trait,
      };
      return el;
    })).toList();

    let endElem: EndOfItemsElement = {
      type: 'END_OF_ITEMS',
      id: 'END_OF_ITEMS',
      numElems: data.size,
    };
    els = els.push(endElem).toList();

    const ds = this.ds.cloneWithRows(els.toJS());
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
    width: Window.WIDTH,
    // TODO: doesn't take into account header height
    height: Window.HEIGHT - Constants.statusBarHeight,
    top: 0,
    left: 0,
    alignSelf: 'stretch',
    justifyContent: 'center',
    backgroundColor: 'white',
  },
  searchItem: {
    flexDirection: 'row',
    alignItems: 'center',
    borderBottomWidth: 0.5,
    borderColor: '#909090',
    padding: 10,
  },
  traitContainer: {
    paddingLeft: 10,
    paddingRight: 20,
  },
  traitText: {
    fontSize: 16,
  },
  typeText: {
    fontSize: 12,
  },
  errorText: {
    color: Colors.HIVE_ERROR,
    paddingTop: 5,
  },
  boldText: {
    fontWeight: '900',
  },
  noMoreResults: {
    marginTop: 10,
    paddingHorizontal: 10,
    alignItems: 'center',
  },
  hint: {
    paddingVertical: 15,
    paddingHorizontal: 10,
    alignItems: 'center',
    borderBottomWidth: 0.5,
    borderColor: '#909090',
  },
});

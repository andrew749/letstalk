import React, { SFC } from 'react';
import { Text, View, ViewStyle, Platform, StyleSheet, TouchableOpacity } from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import Autocomplete from 'react-native-autocomplete-input';
import { FormValidationMessage, FormInputProps, FormLabel } from 'react-native-elements';

export interface DataItem {
  readonly name: string;
  readonly id: number | string;
}

interface SelectItem {
  readonly type: 'ITEM';
  readonly name: string;
  readonly id: number | string;
}

interface SelectCustomItem {
  readonly type: 'CUSTOM_ITEM';
  readonly name: string;
}

export type Select = SelectItem | SelectCustomItem;

type Props = WrappedFieldProps & {
  label: string;
  onQueryChange(query: string, setData: (items: Array<DataItem>) => void): Promise<void>;
  placeholder?: string;
  containerStyle?: ViewStyle | Array<ViewStyle>;
}

interface State {
  readonly query: string;
  readonly showSuggestions: boolean;
  readonly items: Array<DataItem>,
}

const initialState = {
  query: '',
  showSuggestions: false,
  items: [] as Array<DataItem>,
}

// Greedily searches for the earliest characters in name that match some prefix of the query.
// Returns [start, end) of the match if one exists, otherwise null.
export function greedyMatch(name: string, query: string): Array<number> | null {
  if (name.length === 0 || query.length === 0) return null;
  name = name.toLowerCase();
  query = query.toLowerCase();
  let start = 0;
  for (;; start++) {
    if (start >= name.length) return null; // couldn't find character
    if (name.charAt(start) === query.charAt(0)) break;
  }
  let end = start + 1;
  for (;; end++) {
    const queryIdx = end - start;
    if (end >= name.length || queryIdx >= query.length) break;
    if (name.charAt(end) !== query.charAt(queryIdx)) break;
  }
  return [start, end];
}

class AutocompleteInput extends React.Component<Props, State> {
  private autocompleteRef: any;
  constructor(props: Props) {
    super(props);
    this.state = initialState;

    this.onChangeText = this.onChangeText.bind(this);
    this.renderItem = this.renderItem.bind(this);
    this.changeValue = this.changeValue.bind(this);
    this.onEndEditing = this.onEndEditing.bind(this);
  }

  focus() {
    // @ts-ignore
    this.autocompleteRef.current.focus();
  }

  // Decides what value to give to `onChange` depending on the current query and items. If items
  // contains an item with a name that has the same value as query (ignoring case), then we set
  // the value to that item. Otherwise, we set the value to a custom item with the name being the
  // current query.
  private changeValue() {
    const { query, items } = this.state;
    const { onChange } = this.props.input;
    if (query !== '') {
      const found = this.state.items.find(item => {
        return item.name.trim().toLowerCase() === query.trim().toLowerCase()
      });
      const newValue = !found ? { type: 'CUSTOM_ITEM', name: query.trim() } :
        { type: 'ITEM', id: found.id, name: found.name };
      onChange(newValue);
    } else {
      onChange(null);
    }
  }

  private onChangeText(text: string) {
    this.setState({ query: text }, async () => {
      this.changeValue();
      await this.props.onQueryChange(text, items => this.setState({ items }));
      this.changeValue(); // Try again once we have new items
    });
  }

  private renderItem(item: DataItem) {
    const { query } = this.state;
    let text = null;
    const { name, id } = item;
    const matched = greedyMatch(item.name, query);
    if (!!matched) {
      text = (
        <Text>
          {name.substring(0, matched[0])}
          <Text style={{fontWeight: "900"}}>{name.substring(matched[0], matched[1])}</Text>
          {name.substring(matched[1])}
        </Text>
      );
    } else {
      text = <Text>{name}</Text>;
    }
    const onPress = () => {
      this.onChangeText(name);
      this.autocompleteRef.blur();
    };
    return (
      <TouchableOpacity onPress={onPress} style={styles.itemContainer}>
        <Text style={styles.itemText}>{ text }</Text>
      </TouchableOpacity>
    );
  }

  private onEndEditing() {
    const { value } = this.props.input;
    if (value !== null) {
      this.setState({ query: value.name });
    }
  }

  render() {
    const props = this.props;
    const { label, containerStyle } = props;
    const { value, onBlur, onFocus } = this.props.input;
    const { error, touched, warning } = props.meta;
    const {
      items,
      query,
      showSuggestions,
    } = this.state;
    return (
      <View style={[containerStyle]}>
        {label && <FormLabel>{label}</FormLabel>}
        <View style={styles.outerContainer}>
          <View style={styles.autocompleteContainer}>
            <Autocomplete
              ref={(ref: any) => this.autocompleteRef = ref}
              value={query}
              autoCorrect={true}
              data={items}
              renderItem={this.renderItem}
              hideResults={!showSuggestions}
              onEndEditing={this.onEndEditing}
              onFocus={() => {
                this.setState({ showSuggestions: true });
                onFocus(undefined);
              }}
              onBlur={() => {
                this.setState({ showSuggestions: false });
                onBlur(undefined);
              }}
              onChangeText={this.onChangeText}
              {...this.props}
              listStyle={{maxHeight: 200}}
              placeholder={this.props.placeholder}
            />
          </View>
        </View>
        {touched && (
          (error && <FormValidationMessage>{error}</FormValidationMessage>) ||
          (warning && <FormValidationMessage>{warning}</FormValidationMessage>))}
      </View>
    );
  }
};

const styles = StyleSheet.create({
  autocompleteContainer: {
    flex: 1,
    left: 0,
    position: 'absolute',
    right: 0,
    top: 0,
    zIndex: 1,
  },
  outerContainer: {
    margin: 10,
    marginBottom: 5,
    height: 40,
  },
  itemContainer: {
    padding: 10,
  },
  itemText: {
    fontSize: 16,
  }
});

export default AutocompleteInput;

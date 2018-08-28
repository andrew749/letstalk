import React, { SFC } from 'react';
import { Text, View, Platform, StyleSheet, TouchableOpacity } from 'react-native';
import { WrappedFieldProps } from 'redux-form';
import Autocomplete from 'react-native-autocomplete-input';
import { FormValidationMessage, FormInputProps, FormLabel } from 'react-native-elements';

interface DataItem {
  readonly name: string;
  readonly id: number | string;
}

interface CustomItemElement {
  readonly type: 'CUSTOM_ITEM';
}

interface ItemElement extends DataItem {
  readonly type: 'ITEM';
}

type Element = CustomItemElement | ItemElement;

interface SelectItem {
  readonly type: 'ITEM';
  readonly name: string;
  readonly id: number | string;
}

interface SelectCustomItem {
  readonly type: 'CUSTOM_ITEM';
  readonly query: string;
}

export type Select = SelectItem | SelectCustomItem;

type Props = WrappedFieldProps & {
  label: string;
  onQueryChange(query: string): Promise<Array<DataItem>>;
  allowCustom?: boolean;
  placeholder?: string;
}

interface State {
  readonly query: string;
  readonly showSuggestions: boolean;
  readonly data: Array<Element>,
}

const initialState = {
  query: '',
  showSuggestions: false,
  data: [] as Array<Element>,
}

// Greedily searches for the earilest characters in name that match some prefix of the query.
function greedyMatch(name: string, query: string): Array<number> | null {
  if (name.length === 0 || query.length === 0) return null;
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
  }

  focus() {
    // @ts-ignore
    this.autocompleteRef.current.focus();
  }

  private showCustom(query: string, items: Array<ItemElement>): boolean {
    if (query === '' || !this.props.allowCustom) return false;
    const names = items.map(item => item.name.trim().toLowerCase());
    return names.indexOf(query.trim().toLowerCase()) === -1;
  }

  private async onChangeText(text: string) {
    this.setState({ query: text });
    const res = await this.props.onQueryChange(text);
    let items: Array<ItemElement> = res.map(({ id, name }) => {
      const item: ItemElement = { id, name, type: 'ITEM' };
      return item;
    });
    const data: Array<Element> =
      this.showCustom(text, items) ? [{ type: 'CUSTOM_ITEM' } as Element].concat(items) : items;
    this.setState({ data });
  }

  private renderItem(item: Element) {
    const { onChange } = this.props.input;
    const { query } = this.state;
    let text = null;
    let onPress: () => void = null;
    switch (item.type) {
      case 'ITEM':
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
        onPress = () => {
          this.onChangeText(name);
          onChange({ type: 'ITEM', id, name });
          this.autocompleteRef.blur();
        };
        break;
      case 'CUSTOM_ITEM':
        if (query === '') return null; // Race with autocomplete results
        text = <Text>{'Add '}<Text style={{fontWeight: "900"}}>{query}</Text></Text>
        onPress = () => {
          this.onChangeText(query);
          onChange({ type: 'CUSTOM_ITEM', query });
          this.autocompleteRef.blur();
        };
        break;
      default:
        // Ensure exhaustiveness of select
        const _: never = item;
    }
    return (
      <TouchableOpacity onPress={onPress} style={styles.itemContainer}>
        <Text style={styles.itemText}>{ text }</Text>
      </TouchableOpacity>
    );
  }

  render() {
    const props = this.props;
    const { label } = props;
    const { value, onBlur, onFocus } = this.props.input;
    const { error, touched, warning } = props.meta;
    const {
      data,
      query,
      showSuggestions,
    } = this.state;
    return (
      <View style={styles.container}>
        {label && <FormLabel>{label}</FormLabel>}
        <View style={styles.outerContainer}>
          <View style={styles.autocompleteContainer}>
            <Autocomplete
              ref={(ref: any) => this.autocompleteRef = ref}
              value={query}
              autoCorrect={true}
              data={data}
              renderItem={this.renderItem}
              hideResults={!showSuggestions}
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
    zIndex: 1,
  },
  container: {
    zIndex: 1,
  },
  itemContainer: {
    padding: 10,
  },
  itemText: {
    fontSize: 16,
  }
});

export default AutocompleteInput;

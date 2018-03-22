// Taken from https://github.com/jorilallo/react-native-emoji since the npm package was broken.
// Can use the package once the maintainer updates it.
import React from 'react';
import { Text } from 'react-native';
import nodeEmoji from 'node-emoji';

interface Props {
  name: string;
}

const Emoji: React.SFC<Props> = props => {
  const emoji = nodeEmoji.get(props.name);
  return (<Text>{ emoji }</Text>);
};

export default Emoji;

import * as React from 'react';
import App from '../index';

import * as TestRenderer from 'react-test-renderer';

// mock out module method that doesnt exist in testing since no DOM
import {Linking} from 'expo';
Linking.addEventListener = jest.fn();

it('renders without crashing', () => {
  const rendered = TestRenderer.create(<App />).toJSON();
  expect(rendered).toBeTruthy();
})

import * as React from 'react';
import App from '../index';

import * as TestRenderer from 'react-test-renderer';

it('renders without crashing', () => {
  const rendered = TestRenderer.create(<App />).toJSON();
  expect(rendered).toBeTruthy();
})

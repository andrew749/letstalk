// TODO: Actually add types for this
declare module 'react-native-simple-picker';

declare module 'expo-analytics' {
  import {Hit} from 'expo-analytics';

  class Analytics {
    constructor(trackingId: string);
    hit(pageHit: Hit): Promise<void>;
  }

  class Hit {
    // TODO replace with better typings
    constructor(props: any);
  }

  class PageHit extends Hit {}

  class Event extends Hit {
    constructor(category: string, action: string, label: string, value: number): void;
  }
}

import { Analytics, PageHit, Event } from 'expo-analytics';
import {ANALYTICS_ID} from './constants';

export class AnalyticsHelper {
  private analytics: Analytics;

  static _instance: AnalyticsHelper;
  static getInstance(): AnalyticsHelper {
    if (AnalyticsHelper._instance == null) {
      AnalyticsHelper._instance = new AnalyticsHelper();
    }
    return AnalyticsHelper._instance;
  }

  constructor() {
    this.analytics = new Analytics(ANALYTICS_ID);
  }

  recordPage( pageName: string ) {
    this.analytics.hit(new PageHit(pageName))
      .then(() => console.log("Successfully logged page view: " + pageName))
      .catch(e => console.log(e.message));
  }

  recordAction( category: string, action: string, label: string, value: number ) {
    this.analytics.hit(new Event(category, action, label, value))
      .then(() => console.log("Successfully logged action: " + action))
      .catch(e => console.log(e.message));
  }
}

export default AnalyticsHelper;

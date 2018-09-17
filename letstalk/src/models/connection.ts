export enum IntentTypes {
  SEARCH = 'SEARCH',
  REC_GENERAL = 'RECOMMENDATION_GENERAL',
  REC_COHORT = 'RECOMMENDATION_COHORT',
}

export interface Connection {
  readonly userId: number;
  readonly intentType: IntentTypes;
  readonly searchedTrait: string | null;
  readonly message: string | null;
}

export interface ConnectionIntent {
  readonly intentType: IntentTypes;
  readonly searchedTrait: string | null;
}

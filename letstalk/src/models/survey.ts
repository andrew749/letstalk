import Immutable from 'immutable';

export interface SurveyOption {
  readonly key: string
  readonly text: string
}

export interface SurveyQuestion {
  readonly key: string
  readonly prompt: string
  readonly options: Immutable.List<SurveyOption>
}

export type SurveyQuestions =  Immutable.List<SurveyQuestion>
export type SurveyResponses = Immutable.Map<string, string>;

export interface Survey {
  readonly group: string
  readonly version: number
  readonly questions: SurveyQuestions
  readonly responses?: SurveyResponses
}

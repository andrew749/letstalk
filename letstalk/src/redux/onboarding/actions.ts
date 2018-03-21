import { Action } from 'redux'

export enum Step {
  COHORT = 0,
  MY_VECTOR,
  YOUR_VECTOR,
}

export enum TypeKeys {
  SET_STEP = 'ONBOARDING/SET_STEP',
}

export interface SetStepAction extends Action {
  readonly type: TypeKeys.SET_STEP;
  readonly step: Step;
}

export function setStep(step: Step): SetStepAction {
  return {
    type: TypeKeys.SET_STEP,
    step,
  };
}

export type ActionTypes =
  | SetStepAction;

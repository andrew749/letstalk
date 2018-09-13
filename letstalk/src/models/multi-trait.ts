import { CohortV2 } from './cohort';
import { Position } from './position';
import { SimpleTrait } from './simple-trait';

export enum MultiTraitTypes {
  COHORT       = "COHORT",
  POSITION     = "POSITION",
  SIMPLE_TRAIT = "SIMPLE_TRAIT",
}


interface BaseMultiTrait {
  readonly traitName: string;
}

interface CohortMultiTrait extends CohortV2, BaseMultiTrait {
  readonly traitType: MultiTraitTypes.COHORT;
};

interface PositionMultiTrait extends Position, BaseMultiTrait {
  readonly traitType: MultiTraitTypes.POSITION;
}

interface SimpleTraitMultiTrait extends SimpleTrait, BaseMultiTrait {
  readonly traitType: MultiTraitTypes.SIMPLE_TRAIT;
}

export type MultiTrait =
  | CohortMultiTrait
  | PositionMultiTrait
  | SimpleTraitMultiTrait

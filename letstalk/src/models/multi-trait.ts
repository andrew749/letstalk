import { CohortV2 } from './cohort';
import { Position } from './position';
import { SimpleTraitDenormalized } from './simple-trait';
import { Group } from './group';

export enum MultiTraitTypes {
  COHORT       = 'COHORT',
  POSITION     = 'POSITION',
  SIMPLE_TRAIT = 'SIMPLE_TRAIT',
  GROUP        = 'GROUP',
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

interface SimpleTraitMultiTrait extends SimpleTraitDenormalized, BaseMultiTrait {
  readonly traitType: MultiTraitTypes.SIMPLE_TRAIT;
}

interface GroupMultiTrait extends Group, BaseMultiTrait {
  readonly traitType: MultiTraitTypes.GROUP;
}

export type MultiTrait =
  | CohortMultiTrait
  | PositionMultiTrait
  | SimpleTraitMultiTrait
  | GroupMultiTrait

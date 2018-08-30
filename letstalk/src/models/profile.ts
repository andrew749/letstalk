import Immutable from 'immutable';

import { Cohort } from './cohort';
import { UserAdditionalData, UserContactInfo, UserPersonalInfo } from './user';
import { UserPosition } from './position';
import { UserSimpleTrait } from './simple-trait';

export interface ProfileData extends UserAdditionalData, UserPersonalInfo, UserContactInfo, Cohort {
  readonly userPositions: Immutable.List<UserPosition>;
  readonly userSimpleTraits: Immutable.List<UserSimpleTrait>;
}

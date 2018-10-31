type SimpleTraitType = string;

const SIMPLE_TRAIT_TYPE_INTEREST:     SimpleTraitType = "INTEREST";
const SIMPLE_TRAIT_TYPE_EXPERIENCE:   SimpleTraitType = "EXPERIENCE";
const SIMPLE_TRAIT_TYPE_RELIGION:     SimpleTraitType = "RELIGION";
const SIMPLE_TRAIT_TYPE_RACE:         SimpleTraitType = "RACE";
const SIMPLE_TRAIT_TYPE_UNDETERMINED: SimpleTraitType = "UNDETERMINED";

export interface SimpleTrait {
  readonly id: number;
  readonly name: string;
  readonly type: SimpleTraitType;
  readonly isSensitive: boolean;
};

export {
  SIMPLE_TRAIT_TYPE_INTEREST,
  SIMPLE_TRAIT_TYPE_EXPERIENCE,
  SIMPLE_TRAIT_TYPE_RELIGION,
  SIMPLE_TRAIT_TYPE_RACE,
  SIMPLE_TRAIT_TYPE_UNDETERMINED,
};

// Used in other entities such as UserSimpleTrait and SimpleTraitMultiTrait
export interface SimpleTraitDenormalized {
  readonly simpleTraitId: number;
  readonly simpleTraitName: string;
  readonly simpleTraitType: SimpleTraitType;
  readonly simpleTraitIsSensitive: boolean;
}

export interface UserSimpleTrait extends SimpleTraitDenormalized {
  readonly id: number;
}

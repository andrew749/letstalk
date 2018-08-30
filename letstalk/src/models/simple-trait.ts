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

export interface UserSimpleTrait {
  readonly id: number;
  readonly simpleTraitId: number;
  readonly simpleTraitName: number;
  readonly simpleTraitType: SimpleTraitType;
  readonly simpleTraitIsSensitive: boolean;
}

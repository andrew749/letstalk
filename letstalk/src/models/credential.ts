import Immutable from 'immutable';

export enum CredentialPositionType {
  COOP = 0,
  CLUB,
  SPORT,
  COHORT,
};

export interface CredentialPosition {
  readonly id: number;
  readonly name: string;
  readonly type: CredentialPositionType;
}

export enum CredentialOrganizationType {
  COOP = 0,
  CLUB,
  SPORT,
  COHORT,
}

export interface CredentialFilterableElement {
  readonly id: [number, number]; // [positionId, organizationId]
  readonly value: string;
}

export interface Credential {
  readonly positionId: number;
  readonly positionName: string;
  readonly organizationId: number;
  readonly organizationName: string;
}

export interface CredentialOrganization {
  readonly id: number;
  readonly name: string;
  readonly type: CredentialOrganizationType;
}

export interface ValidCredentialPair {
  readonly positionType: CredentialPositionType;
  readonly organizationType: CredentialOrganizationType;
}

export interface CredentialOptions {
  readonly validPairs: Immutable.List<ValidCredentialPair>;
  readonly organizations: Immutable.List<CredentialOrganization>;
  readonly positions: Immutable.List<CredentialPosition>;
}

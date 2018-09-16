export interface Role {
  readonly id: number;
  readonly name: string;
};

type OrganizationType = string;

const	ORGANIZATION_TYPE_COMPANY:      OrganizationType = "COMPANY";
const	ORGANIZATION_TYPE_CLUB:         OrganizationType = "CLUB";
const	ORGANIZATION_TYPE_SPORTS_TEAM:  OrganizationType = "SPORTS_TEAM";
const	ORGANIZATION_TYPE_UNDETERMINED: OrganizationType = "UNDETERMINED";

export interface Organization {
  readonly id: number;
  readonly name: string;
  readonly type: OrganizationType;
};

export {
  ORGANIZATION_TYPE_COMPANY,
  ORGANIZATION_TYPE_CLUB,
  ORGANIZATION_TYPE_SPORTS_TEAM,
  ORGANIZATION_TYPE_UNDETERMINED,
};

export interface Position {
  readonly roleId: number;
  readonly roleName: string;
  readonly organizationId: number;
  readonly organizationName: string;
  readonly organizationType: OrganizationType,
}

export interface UserPosition extends Position {
  readonly id: number;
  readonly startDate: string;
  readonly endDate: string | null;
};

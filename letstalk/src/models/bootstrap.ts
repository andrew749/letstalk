import Immutable from 'immutable';

import { Cohort } from './cohort';
import { UserPersonalInfo } from './user';
import { OnboardingStatus } from './onboarding';
import { Connection } from './connection';

export interface Relationship {
  readonly userId: number;
  readonly userType: number;
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
  readonly fbId: string | null;
  readonly fbLink: string | null;
  readonly phoneNumber: string | null;
  readonly cohort: Cohort | null;
  readonly description: string | null;
  readonly matchingState: MatchingState;
}

export function getHumanReadableUserType(userType: number): string {
  switch (userType) {
    case 1:
      return 'Mentor';
    case 2:
      return 'Mentee';
    case 3:
      return 'Asker';
    case 4:
      return 'Answerer';
  }
  return 'Unknown';
}

export enum MatchingState {
  Unknown = 0,
  Unverified,
  Verified,
  Expired,
}

export type UserState =
  | 'account_created'
  | 'account_email_verified'
  | 'account_setup'
  | 'account_matched';

export interface ConnectionRequestWithName extends Connection {
  readonly firstName: string;
  readonly lastName: string;
}

export interface BootstrapConnection {
  readonly userProfile: Relationship;
  readonly request: Connection;
}

interface BootstrapConnections {
  readonly outgoingRequests: Immutable.List<ConnectionRequestWithName>;
  readonly incomingRequests: Immutable.List<ConnectionRequestWithName>;
  readonly mentors: Immutable.List<BootstrapConnection>;
  readonly mentees: Immutable.List<BootstrapConnection>;
  readonly peers: Immutable.List<BootstrapConnection>;
}

export interface BootstrapData {
  readonly state: UserState;
  readonly cohort: Cohort;
  readonly onboardingStatus: OnboardingStatus;
  readonly connections: BootstrapConnections;
};

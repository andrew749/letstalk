import Immutable from 'immutable';

export interface Relationship {
  readonly userId: number;
  readonly userType: 'mentor' | 'mentee';
  readonly firstName: string;
  readonly lastName: string;
  readonly email: string;
}

export interface Cohort {
  readonly cohortId: number;
  readonly programId: string;
  readonly gradYear: number;
  readonly sequence: string;
}

export interface BootstrapData {
  readonly relationships: Immutable.List<Relationship>;
  readonly state: 'account_created' | 'account_setup' | 'account_matched';
  readonly cohort: Cohort;
};

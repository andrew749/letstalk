import Immutable from 'immutable';

export interface Credential {
  readonly id: number;
  readonly name: string;
}

export type Credentials = Immutable.List<Credential>;

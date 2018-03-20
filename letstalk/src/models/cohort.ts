import Immutable from 'immutable';

const PROGRAMS: Immutable.Map<string, string> = Immutable.Map({
  softwareEngineering: 'Software Engineering',
  computerEngineering: 'Computer Engineering',
});

const SEQUENCES: Immutable.Map<string, string> = Immutable.Map({
  stream4: 'Stream 4',
  stream8: 'Stream 8',
});

const COHORTS: Immutable.List<Immutable.List<any>> = Immutable.fromJS([
  [1, 'softwareEngineering', 'stream8', '2019'],
  [2, 'computerEngineering', 'stream8', '2019'],
  [3, 'computerEngineering', 'stream4', '2019'],
]);

export interface ValueLabel {
  value: string;
  label: string;
}

export function programOptions(): Immutable.List<ValueLabel> {
  return PROGRAMS.map((label, value) => ({ value, label })).toList();
}

export function sequenceOptions(programId: string | null): Immutable.List<ValueLabel> {
  const cohorts = COHORTS.filter(row => !programId || programId === row.get(1));
  return cohorts.map(row => row.get(2)).toSet().map(sequenceId => {
    return { value: sequenceId, label: SEQUENCES.get(sequenceId) };
  }).toList();
}

export function gradYearOptions(programId: string | null, sequenceId: string | null)
  : Immutable.List<ValueLabel> {
  const preCohorts = COHORTS.filter(row => !programId || programId === row.get(1));
  const cohorts = preCohorts.filter(row => !sequenceId || sequenceId === row.get(2));
  return cohorts.map(row => row.get(3)).toSet().map(gradYear => {
    return { value: gradYear, label: gradYear };
  }).toList();
}

export function getCohortId(programId: string , sequenceId: string, gradYear: string): number {
  const row = COHORTS.find(row => {
    return row.get(1) === programId &&
      row.get(2) === sequenceId &&
      row.get(3) === gradYear;
  });
  if (row === null) return null;
  return row.get(0);
}

export interface Cohort {
  readonly cohortId: number;
  readonly programId: string;
  readonly gradYear: number;
  readonly sequence: string;
}

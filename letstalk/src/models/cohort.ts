import Immutable from 'immutable';

const PROGRAMS: Immutable.Map<string, string> = Immutable.Map({
  'SOFTWARE_ENGINEERING': 'Software Engineering',
  'COMPUTER_ENGINEERING': 'Computer Engineering',
});

export function programById(programId: string): string {
  return PROGRAMS.get(programId, programId);
}

const SEQUENCES: Immutable.Map<string, string> = Immutable.Map({
  '4STREAM': '4 Stream',
  '8STREAM': '8 Stream',
});

export function sequenceById(sequenceId: string): string {
  return SEQUENCES.get(sequenceId, sequenceId);
}

export interface ValueLabel {
  readonly value: any;
  readonly label: string;
}

type ValueLabels = Immutable.List<ValueLabel>;

export interface Cohort {
  readonly cohortId: number;
  readonly programId: string;
  readonly sequenceId: string;
  readonly gradYear: number;
}

type Cohorts = Immutable.List<Cohort>;

export function programOptions(cohorts: Cohorts): ValueLabels {
  return cohorts.map(row => row.programId).toSet().map(
    programId => ({ value: programId, label: PROGRAMS.get(programId) })
  ).toList();
}

function filteredCohorts(
  cohorts: Cohorts,
  programId?: string,
  sequenceId?: string,
): Cohorts {
  return cohorts.filter(row => {
    return (!programId || programId === row.programId) &&
      (!sequenceId || sequenceId === row.sequenceId)
  }).toList();
}

export function sequenceOptions(
  cohorts: Cohorts,
  programId: string,
): ValueLabels {
  return filteredCohorts(cohorts, programId).map(
    row => row.sequenceId
  ).toSet().map(
    sequenceId => ({ value: sequenceId, label: SEQUENCES.get(sequenceId) })
  ).toList();
}

export function gradYearOptions(
  cohorts: Cohorts,
  programId: string,
  sequenceId: string,
): ValueLabels {
  return filteredCohorts(cohorts, programId, sequenceId).map(
    row => row.gradYear
  ).toSet().map(gradYear => {
    const gradYearStr = String(gradYear);
    return { value: gradYear, label: gradYearStr };
  }).toList();
}

export function getCohortId(
  cohorts: Cohorts,
  programId: string,
  sequenceId: string,
  gradYear: number,
): number {
  const row = cohorts.find(row => {
    return row.programId === programId &&
      row.sequenceId === sequenceId &&
      row.gradYear === gradYear;
  });
  if (row === null) return null;
  return row.cohortId;
}

var cohortId = 1;
const COHORTS = Immutable.List(
  [2018, 2019, 2020, 2021, 2022, 2023],
).flatMap(gradYear => {
  const cohorts = Immutable.List([
    { cohortId: cohortId, programId: 'SOFTWARE_ENGINEERING', sequenceId: '8STREAM', gradYear},
    { cohortId: cohortId + 1, programId: 'COMPUTER_ENGINEERING', sequenceId: '8STREAM', gradYear},
    { cohortId: cohortId + 2, programId: 'COMPUTER_ENGINEERING', sequenceId: '4STREAM', gradYear},
  ]);
  cohortId = cohortId + 3;
  return cohorts;
}).toList();

export { COHORTS };

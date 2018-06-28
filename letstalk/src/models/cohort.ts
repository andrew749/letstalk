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

// TODO: fetch cohorts from server and store in redux state
const COHORTS = Immutable.List<Cohort>([
  { cohortId:  4, programId: 'SOFTWARE_ENGINEERING', gradYear: 2018, sequenceId: '8STREAM'},
  { cohortId:  5, programId: 'COMPUTER_ENGINEERING', gradYear: 2018, sequenceId: '8STREAM'},
  { cohortId:  6, programId: 'COMPUTER_ENGINEERING', gradYear: 2018, sequenceId: '4STREAM'},
  { cohortId:  1, programId: 'SOFTWARE_ENGINEERING', gradYear: 2019, sequenceId: '8STREAM'},
  { cohortId:  2, programId: 'COMPUTER_ENGINEERING', gradYear: 2019, sequenceId: '8STREAM'},
  { cohortId:  3, programId: 'COMPUTER_ENGINEERING', gradYear: 2019, sequenceId: '4STREAM'},
  { cohortId:  7, programId: 'SOFTWARE_ENGINEERING', gradYear: 2020, sequenceId: '8STREAM'},
  { cohortId:  8, programId: 'COMPUTER_ENGINEERING', gradYear: 2020, sequenceId: '8STREAM'},
  { cohortId:  9, programId: 'COMPUTER_ENGINEERING', gradYear: 2020, sequenceId: '4STREAM'},
  { cohortId: 10, programId: 'SOFTWARE_ENGINEERING', gradYear: 2021, sequenceId: '8STREAM'},
  { cohortId: 11, programId: 'COMPUTER_ENGINEERING', gradYear: 2021, sequenceId: '8STREAM'},
  { cohortId: 12, programId: 'COMPUTER_ENGINEERING', gradYear: 2021, sequenceId: '4STREAM'},
  { cohortId: 13, programId: 'SOFTWARE_ENGINEERING', gradYear: 2022, sequenceId: '8STREAM'},
  { cohortId: 14, programId: 'COMPUTER_ENGINEERING', gradYear: 2022, sequenceId: '8STREAM'},
  { cohortId: 15, programId: 'COMPUTER_ENGINEERING', gradYear: 2022, sequenceId: '4STREAM'},
  { cohortId: 16, programId: 'SOFTWARE_ENGINEERING', gradYear: 2023, sequenceId: '8STREAM'},
  { cohortId: 17, programId: 'COMPUTER_ENGINEERING', gradYear: 2023, sequenceId: '8STREAM'},
  { cohortId: 18, programId: 'COMPUTER_ENGINEERING', gradYear: 2023, sequenceId: '4STREAM'},
]);

export { COHORTS };

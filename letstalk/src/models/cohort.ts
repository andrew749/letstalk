import Immutable from 'immutable';

const PROGRAMS: Immutable.Map<string, string> = Immutable.Map({
  'SOFTWARE_ENGINEERING': 'Software Engineering',
  'COMPUTER_ENGINEERING': 'Computer Engineering',
});

const SEQUENCES: Immutable.Map<string, string> = Immutable.Map({
  '4STREAM': '4 Stream',
  '8STREAM': '8 Stream',
});

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

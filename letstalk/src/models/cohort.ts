import Immutable from 'immutable';

const PROGRAMS: Immutable.Map<string, string> = Immutable.Map({
  'SOFTWARE_ENGINEERING': 'Software Engineering',
  'COMPUTER_ENGINEERING': 'Computer Engineering',
  'ARCHITECTURAL_ENGINEERING': 'Architecture Engineering',
  'BIOMEDICAL_ENGINEERING': 'Biomedical Engineering',
  'CHEMICAL_ENGINEERING': 'Chemical Engineering',
  'CIVIL_ENGINEERING': 'Civil Engineering',
  'ELECTRICAL_ENGINEERING': 'Electrical Engineering',
  'ENVIRONMENTAL_ENGINEERING': 'Environmental Engineering',
  'GEOLOGICAL_ENGINEERING': 'Geological Engineering',
  'MANAGEMENT_ENGINEERING': 'Management Engineering',
  'MECHANICAL_ENGINEERING': 'Mechanical Engineering',
  'MECHATRONICS_ENGINEERING': 'Mechatronics Engineering',
  'NANOTECHNOLOGY_ENGINEERING': 'Nanotechnology Engineering',
  'SYSTEMS_DESIGN_ENGINEERING': 'Systems Design Engineering',
  'ACCOUNTING_AND_FINANCIAL_MANAGEMENT': 'Accounting and Financial Management',
  'ACTUARIAL_SCIENCE': 'Actuarial Science',
  'ANTHROPOLOGY': 'Anthropology',
  'APPLIED_MATHEMATICS': 'Applied Mathematics',
  'ARCHITECTURE': 'Architecture',
  'BIOCHEMISTRY': 'Biochemistry',
  'BIOLOGY': 'Biology',
  'BIOMEDICAL_SCIENCES': 'Biomedical Sciences',
  'BIOSTATISTICS': 'Biostatistics',
  'BIOTECHNOLOGY/CHARTERED_PROFESSIONAL_ACCOUNTANCY': 'Biotechnology/Chartered Professional Accountancy',
  'BIOTECHNOLOGY/ECONOMICS': 'Biotechnology/Economics',
  'BUSINESS_ADMINISTRATION_AND_COMPUTER_SCIENCE_DOUBLE_DEGREE': 'Business Administration and Computer Science Double Degree',
  'BUSINESS_ADMINISTRATION_AND_MATHEMATICS_DOUBLE_DEGREE': 'Business Administration and Mathematics Double Degree',
  'CHEMISTRY': 'Chemistry',
  'CLASSICAL_STUDIES': 'Classical Studies',
  'COMBINATORICS_AND_OPTIMIZATION': 'Combinatorics and Optimization',
  'COMPUTATIONAL_MATHEMATICS': 'Computational Mathematics',
  'COMPUTER_SCIENCE': 'Computer Science',
  'COMPUTING_AND_FINANCIAL_MANAGEMENT': 'Computing and Financial Management',
  'DATA_SCIENCE': 'Data Science',
  'EARTH_SCIENCES': 'Earth Sciences',
  'ECONOMICS': 'Economics',
  'ENGLISH': 'English',
  'ENVIRONMENT_AND_BUSINESS': 'Environment and Business',
  'ENVIRONMENT,_RESOURCES_AND_SUSTAINABILITY': 'Environment, Resources and Sustainability',
  'ENVIRONMENTAL_SCIENCE': 'Environmental Science',
  'FINE_ARTS': 'Fine Arts',
  'FRENCH': 'French',
  'GENDER_AND_SOCIAL_JUSTICE': 'Gender and Social Justice',
  'GEOGRAPHY_AND_AVIATION': 'Geography and Aviation',
  'GEOGRAPHY_AND_ENVIRONMENTAL_MANAGEMENT': 'Geography and Environmental Management',
  'GEOMATICS': 'Geomatics',
  'GERMAN': 'German',
  'GLOBAL_BUSINESS_AND_DIGITAL_ARTS': 'Global Business and Digital Arts',
  'HEALTH_STUDIES': 'Health Studies',
  'HISTORY': 'History',
  'HONOURS_ARTS': 'Honours Arts',
  'HONOURS_ARTS_AND_BUSINESS': 'Honours Arts and Business',
  'HONOURS_SCIENCE': 'Honours Science',
  'INFORMATION_TECHNOLOGY_MANAGEMENT': 'Information Technology Management',
  'INTERNATIONAL_DEVELOPMENT': 'International Development',
  'KINESIOLOGY': 'Kinesiology',
  'KNOWLEDGE_INTEGRATION': 'Knowledge Integration',
  'LEGAL_STUDIES': 'Legal Studies',
  'LIBERAL_STUDIES': 'Liberal Studies',
  'LIFE_PHYSICS': 'Life Physics',
  'LIFE_SCIENCES': 'Life Sciences',
  'MATERIALS_AND_NANOSCIENCES': 'Materials and Nanosciences',
  'MATHEMATICAL_ECONOMICS': 'Mathematical Economics',
  'MATHEMATICAL_FINANCE': 'Mathematical Finance',
  'MATHEMATICAL_OPTIMIZATION': 'Mathematical Optimization',
  'MATHEMATICAL_PHYSICS': 'Mathematical Physics',
  'MATHEMATICAL_STUDIES': 'Mathematical Studies',
  'MATHEMATICS': 'Mathematics',
  'MATHEMATICS/BUSINESS_ADMINISTRATION': 'Mathematics/Business Administration',
  'MATHEMATICS/CHARTERED_PROFESSIONAL_ACCOUNTANCY': 'Mathematics/Chartered Professional Accountancy',
  'MATHEMATICS/FINANCIAL_ANALYSIS_AND_RISK_MANAGEMENT': 'Mathematics/Financial Analysis and Risk Management',
  'MEDICINAL_CHEMISTRY': 'Medicinal Chemistry',
  'MEDIEVAL_STUDIES': 'Medieval Studies',
  'MUSIC': 'Music',
  'PEACE_AND_CONFLICT_STUDIES': 'Peace and Conflict Studies',
  'PHILOSOPHY': 'Philosophy',
  'PHYSICAL_SCIENCES': 'Physical Sciences',
  'PHYSICS': 'Physics',
  'PHYSICS_AND_ASTRONOMY': 'Physics and Astronomy',
  'PLANNING': 'Planning',
  'POLITICAL_SCIENCE': 'Political Science',
  'PSYCHOLOGY': 'Psychology',
  'PUBLIC_HEALTH': 'Public Health',
  'PURE_MATHEMATICS': 'Pure Mathematics',
  'RECREATION_AND_LEISURE_STUDIES': 'Recreation and Leisure Studies',
  'RECREATION_AND_SPORT_BUSINESS': 'Recreation and Sport Business',
  'SCIENCE_AND_AVIATION': 'Science and Aviation',
  'SCIENCE_AND_BUSINESS': 'Science and Business',
  'SPANISH': 'Spanish',
  'SPEECH_COMMUNICATION': 'Speech Communication',
  'STATISTICS': 'Statistics',
  'THEATRE_AND_PERFORMANCE': 'Theatre and Performance',
  'THERAPEUTIC_RECREATION': 'Therapeutic Recreation',
  'TOURISM_DEVELOPMENT': 'Tourism Development',
  'OTHER': 'Other',
  'ALUM': 'Alum',
});

export function programById(programId: string): string {
  return PROGRAMS.get(programId, programId);
}

const SEQUENCES: Immutable.Map<string, string> = Immutable.Map({
  '4STREAM': '4 Stream',
  '8STREAM': '8 Stream',
  'OTHER': 'Other',
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

// TODO: Migrate all parts of code to V2
export interface CohortV2 {
  readonly cohortId: number;
  readonly programId: string;
  readonly programName: string;
  readonly gradYear: number;
  readonly isCoop: boolean;
  readonly sequenceId: string | null;
  readonly sequenceName: string | null;
}

type Cohorts = Immutable.List<Cohort>;

export function programOptions(cohorts: Cohorts): ValueLabels {
  return cohorts.map(row => row.programId).toSet().map(
    programId => ({ value: programId, label: PROGRAMS.get(programId) })
  ).sortBy(program => program.label).toList();
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
  ).sortBy(sequence => sequence.label).toList();
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
  }).sortBy(gradYear => gradYear.value).toList();
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

export function humanReadableCohort(cohort: CohortV2): string {
  const {
    programName,
    sequenceName,
    sequenceId,
    gradYear,
  } = cohort;
  let cohortText = programName + ' ' + gradYear;
  if (!!sequenceId && !!sequenceName && sequenceId !== 'OTHER') {
    cohortText = cohortText + ' ' + sequenceName;
  }
  return cohortText;
}

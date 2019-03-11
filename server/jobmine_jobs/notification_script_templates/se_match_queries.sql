CREATE VIEW `se_matches` as select user_one_id as mentor_id, user_two_id as mentee_id from mentorships inner join connections on connections.connection_id = mentorships.connection_id where mentorships.created_at > '2019-02-21 02:00:53';
CREATE VIEW `se_mentors` as select distinct users.user_id, users.first_name, users.last_name, cohorts.grad_year from `se_matches` inner join user_cohorts on user_cohorts.user_id = `se_matches`.mentor_id inner join cohorts on cohorts.cohort_id = user_cohorts.cohort_id inner join users on users.user_id = `se_matches`.mentor_id where cohorts.program_id = 'SOFTWARE_ENGINEERING';
CREATE VIEW `se_mentees` as select distinct users.user_id, users.first_name, users.last_name, cohorts.grad_year from `se_matches` inner join user_cohorts on user_cohorts.user_id = `se_matches`.mentee_id inner join cohorts on cohorts.cohort_id = user_cohorts.cohort_id inner join users on users.user_id = `se_matches`.mentee_id where cohorts.program_id = 'SOFTWARE_ENGINEERING';
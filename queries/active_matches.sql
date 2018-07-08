select distinct
  users.first_name,
  name,
  mentors.first_name
from users
inner join request_matchings on asker=user_id
inner join credentials on credential_id=credentials.id
inner join users mentors on answerer=mentors.user_id
where request_matchings.deleted_at IS NULL;

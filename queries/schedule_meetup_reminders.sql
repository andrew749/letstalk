
-- Connections made between (min_time, max_time).
SET @min_time = TIMESTAMP('2019-01-29 00:00:00');
SET @max_time = TIMESTAMP('2019-01-30 00:00:00');
-- Reminders scheduled randomly between now and (now+schedule_interval)
SET @schedule_interval = 60*60*24*1; -- seconds
INSERT INTO meetup_reminders
    (created_at, updated_at, user_id, match_user_id, type, state, scheduled_at)
  SELECT NOW(),
         NOW(),
         connections.user_one_id,
         connections.user_two_id,
         'INITIAL_MEETING',
         'SCHEDULED',
         ADDDATE(NOW(), INTERVAL FLOOR(RAND()*@schedule_interval) SECOND)
  FROM connections WHERE connections.created_at > @min_time AND connections.created_at < @max_time;
INSERT INTO meetup_reminders
    (created_at, updated_at, user_id, match_user_id, type, state, scheduled_at)
  SELECT NOW(),
         NOW(),
         connections.user_two_id,
         connections.user_one_id,
         'INITIAL_MEETING',
         'SCHEDULED',
         ADDDATE(NOW(), INTERVAL FLOOR(RAND()*@schedule_interval) SECOND)
  FROM connections WHERE connections.created_at > @min_time AND connections.created_at < @max_time;

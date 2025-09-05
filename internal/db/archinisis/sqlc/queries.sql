-- name: GetSporttiIDs :many
SELECT sportti_id
FROM sportti_id_list
ORDER BY sportti_id;

-- name: UpsertAthlete :exec
INSERT INTO athlete (
  national_id, first_name, last_name, initials, date_of_birth, height, weight
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (national_id) DO UPDATE SET
  first_name    = EXCLUDED.first_name,
  last_name     = EXCLUDED.last_name,
  initials      = EXCLUDED.initials,
  date_of_birth = EXCLUDED.date_of_birth,
  height        = EXCLUDED.height,
  weight        = EXCLUDED.weight;

-- name: UpsertMeasurement :exec
INSERT INTO measurement (
  measurement_group_id, measurement_id, national_id, discipline, session_name,
  place, race_id, start_time, stop_time, nb_segments, comment
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9, $10, $11
)
ON CONFLICT (measurement_group_id) DO UPDATE SET
  measurement_id = EXCLUDED.measurement_id,
  national_id    = EXCLUDED.national_id,
  discipline     = EXCLUDED.discipline,
  session_name   = EXCLUDED.session_name,
  place          = EXCLUDED.place,
  race_id        = EXCLUDED.race_id,
  start_time     = EXCLUDED.start_time,
  stop_time      = EXCLUDED.stop_time,
  nb_segments    = EXCLUDED.nb_segments,
  comment        = EXCLUDED.comment;

-- name: UpsertRaceReport :exec
INSERT INTO report (sportti_id, session_id, race_report)
VALUES ($1, $2, $3)
ON CONFLICT (session_id) DO UPDATE SET
  sportti_id  = EXCLUDED.sportti_id,
  race_report = EXCLUDED.race_report;

-- name: GetRaceReportSessionIDsBySporttiID :many
SELECT session_id
FROM report
WHERE sportti_id = $1
ORDER BY session_id DESC;

-- name: GetRaceReport :one
SELECT race_report
FROM report
WHERE sportti_id = $1 AND session_id = $2;
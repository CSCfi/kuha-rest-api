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

-- name: UpsertReport :exec
INSERT INTO report (session_id, race_report)
VALUES ($1, $2)
ON CONFLICT (session_id) DO UPDATE SET
  race_report = EXCLUDED.race_report;


-- name: UpsertReportUser :exec
INSERT INTO report_user (session_id, sportti_id)
VALUES ($1, $2)
ON CONFLICT (session_id, sportti_id) DO NOTHING;

-- name: GetRaceReportSessionIDsBySporttiID :many
SELECT ru.session_id
FROM report_user ru
WHERE ru.sportti_id = $1
ORDER BY ru.session_id DESC;

-- name: GetRaceReport :one
SELECT r.race_report
FROM report r
JOIN report_user ru ON ru.session_id = r.session_id
WHERE ru.sportti_id = $1
  AND r.session_id = $2;

-- name: GetSporttiIDsBySessionID :many
SELECT sportti_id
FROM report_user
WHERE session_id = $1;

-- name: GetAthleteBySporttiID :one
SELECT
  national_id,
  first_name,
  last_name,
  initials,
  date_of_birth,
  height,
  weight
FROM athlete
WHERE national_id = $1;

-- name: GetMeasurementsBySporttiID :many
SELECT
  measurement_group_id,
  measurement_id,
  national_id,
  discipline,
  session_name,
  place,
  race_id,
  start_time,
  stop_time,
  nb_segments,
  comment
FROM measurement
WHERE national_id = $1
ORDER BY measurement_group_id ASC;

-- name: DeleteAthleteByNationalID :one
DELETE FROM athlete
WHERE national_id = $1
RETURNING national_id;

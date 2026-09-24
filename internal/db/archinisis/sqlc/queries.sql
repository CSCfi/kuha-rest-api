-- Athlete

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

-- name: DeleteAthleteByNationalID :one
DELETE FROM athlete
WHERE national_id = $1
RETURNING national_id;

-- Measurement group

-- name: EnsureMeasurementGroup :exec
INSERT INTO measurement_group (measurement_group_id)
VALUES ($1)
ON CONFLICT (measurement_group_id) DO NOTHING;

-- Measurement

-- name: UpsertMeasurement :exec
INSERT INTO measurement (
  measurement_group_id, measurement_id, national_id, discipline, session_name,
  place, race_id, start_time, stop_time, nb_segments, comment
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9, $10, $11
)
ON CONFLICT (measurement_id) DO UPDATE SET
  measurement_group_id = EXCLUDED.measurement_group_id,
  national_id          = EXCLUDED.national_id,
  discipline           = EXCLUDED.discipline,
  session_name         = EXCLUDED.session_name,
  place                = EXCLUDED.place,
  race_id              = EXCLUDED.race_id,
  start_time           = EXCLUDED.start_time,
  stop_time            = EXCLUDED.stop_time,
  nb_segments          = EXCLUDED.nb_segments,
  comment              = EXCLUDED.comment;

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
ORDER BY measurement_group_id ASC, measurement_id ASC;

-- name: GetMeasurementByMeasurementID :one
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
WHERE measurement_id = $1;

-- name: DeleteMeasurementByMeasurementID :one
DELETE FROM measurement
WHERE measurement_id = $1
RETURNING national_id;

-- Report (one row per sportti_id + session_id; no unique constraint in the DB,
-- so the upsert is done as UPDATE-then-INSERT inside a transaction)

-- name: UpdateReport :execrows
UPDATE report
SET race_report = $3
WHERE sportti_id = $1
  AND session_id = $2;

-- name: InsertReport :exec
INSERT INTO report (sportti_id, session_id, race_report)
VALUES ($1, $2, $3);

-- name: GetRaceReportSessionIDsBySporttiID :many
SELECT DISTINCT session_id
FROM report
WHERE sportti_id = $1
  AND session_id IS NOT NULL
ORDER BY session_id DESC;

-- name: GetRaceReport :one
SELECT race_report
FROM report
WHERE sportti_id = $1
  AND session_id = $2
ORDER BY report_id DESC
LIMIT 1;

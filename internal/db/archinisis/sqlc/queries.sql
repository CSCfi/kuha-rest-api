-- name: GetSporttiIDs :many
SELECT sportti_id
FROM sportti_id_list
ORDER BY sportti_id;

-- name: UpsertAthlete :one
INSERT INTO athlete (
  national_id, first_name, last_name, initials, date_of_birth, height, weight
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (national_id) DO UPDATE SET
  first_name    = COALESCE(EXCLUDED.first_name,    athlete.first_name),
  last_name     = COALESCE(EXCLUDED.last_name,     athlete.last_name),
  initials      = COALESCE(EXCLUDED.initials,      athlete.initials),
  date_of_birth = COALESCE(EXCLUDED.date_of_birth, athlete.date_of_birth),
  height        = COALESCE(EXCLUDED.height,        athlete.height),
  weight        = COALESCE(EXCLUDED.weight,        athlete.weight)
RETURNING national_id, first_name, last_name, initials, date_of_birth, height, weight;

-- name: UpsertMeasurement :one
INSERT INTO measurement (
  measurement_group_id, measurement_id, national_id, discipline, session_name,
  place, race_id, start_time, stop_time, nb_segments, comment
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9, $10, $11
)
ON CONFLICT (measurement_group_id) DO UPDATE SET
  measurement_id = COALESCE(EXCLUDED.measurement_id, measurement.measurement_id),
  national_id    = COALESCE(EXCLUDED.national_id,    measurement.national_id),
  discipline     = COALESCE(EXCLUDED.discipline,     measurement.discipline),
  session_name   = COALESCE(EXCLUDED.session_name,   measurement.session_name),
  place          = COALESCE(EXCLUDED.place,          measurement.place),
  race_id        = COALESCE(EXCLUDED.race_id,        measurement.race_id),
  start_time     = COALESCE(EXCLUDED.start_time,     measurement.start_time),
  stop_time      = COALESCE(EXCLUDED.stop_time,      measurement.stop_time),
  nb_segments    = COALESCE(EXCLUDED.nb_segments,    measurement.nb_segments),
  comment        = COALESCE(EXCLUDED.comment,        measurement.comment)
RETURNING measurement_group_id, measurement_id, national_id, discipline, session_name,
          place, race_id, start_time, stop_time, nb_segments, comment;

-- name: UpsertRaceReport :one
INSERT INTO report (sportti_id, session_id, race_report)
VALUES ($1, $2, $3)
ON CONFLICT (session_id) DO UPDATE SET
  sportti_id  = COALESCE(EXCLUDED.sportti_id,  report.sportti_id),
  race_report = COALESCE(EXCLUDED.race_report, report.race_report)
RETURNING report_id, sportti_id, session_id, race_report;

-- name: GetRaceReportSessionIDsBySporttiID :many
SELECT session_id
FROM report
WHERE sportti_id = $1
ORDER BY session_id DESC;

-- name: GetRaceReport :one
SELECT race_report
FROM report
WHERE sportti_id = $1 AND session_id = $2;
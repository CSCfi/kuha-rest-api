-- name: GetSporttiIDs :many
SELECT sportti_id FROM sportti_id_list;

-- name: InsertAthlete :exec
INSERT INTO athlete (
  national_id, first_name, last_name, initials,
  date_of_birth, height, weight
) VALUES (
  @national_id, @first_name, @last_name, @initials,
  @date_of_birth, @height, @weight
)
ON CONFLICT (national_id) DO NOTHING;

-- name: InsertMeasurement :exec
INSERT INTO measurement (
  measurement_group_id, measurement_id, national_id,
  discipline, session_name, place, race_id,
  start_time, stop_time, nb_segments, comment
) VALUES (
  @measurement_group_id, @measurement_id, @national_id,
  @discipline, @session_name, @place, @race_id,
  @start_time, @stop_time, @nb_segments, @comment
)
ON CONFLICT (measurement_group_id) DO NOTHING;

-- name: InsertRaceReport :exec
INSERT INTO report (
  sportti_id, session_id, race_report
) VALUES (
  @sportti_id, @session_id, @race_report
)
ON CONFLICT (report_id) DO NOTHING;

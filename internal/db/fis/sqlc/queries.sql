-- name: GetAthletesBySector :many
SELECT Firstname, Lastname, Fiscode
FROM A_competitor
WHERE SectorCode = $1
ORDER BY Fiscode;


-- name: GetNationsBySector :many
SELECT DISTINCT NationCode
FROM A_competitor
WHERE SectorCode = $1
ORDER BY NationCode ASC;


-- name: GetSkiJumpingSeasons :many
SELECT DISTINCT SeasonCode
FROM A_raceJP
ORDER BY SeasonCode DESC;


-- name: GetNordicCombinedSeasons :many
SELECT DISTINCT SeasonCode
FROM A_raceNK
ORDER BY SeasonCode DESC;


-- name: GetSkiJumpingDisciplines :many
SELECT DISTINCT DisciplineCode
FROM A_raceJP
ORDER BY DisciplineCode ASC;


-- name: GetNordicCombinedDisciplines :many
SELECT DISTINCT DisciplineCode
FROM A_raceNK
ORDER BY DisciplineCode ASC;


-- name: GetSkiJumpingCategories :many
SELECT DISTINCT CatCode
FROM A_raceJP
ORDER BY CatCode ASC;


-- name: GetNordicCombinedCategories :many
SELECT DISTINCT CatCode
FROM A_raceNK
ORDER BY CatCode ASC;

-- name: GetCompetitorIDByFiscodeNK :one
SELECT CompetitorID
FROM A_competitor
WHERE Fiscode = $1 AND SectorCode = 'NK';


-- name: GetAthleteResultsNK :many
SELECT 
    A_resultNK.RecID,
    A_resultNK.RaceID,
    A_resultNK.Position,
    A_raceNK.RaceDate,
    A_raceNK.SeasonCode,
    A_raceNK.Distance,
    A_raceNK.Hill,
    A_raceNK.DisciplineCode,
    A_raceNK.CatCode,
    A_raceNK.Place,
    A_resultNK.PosR1,
    A_resultNK.SpeedR1,
    A_resultNK.DistR1,
    A_resultNK.JudPtsR1,
    A_resultNK.WindR1,
    A_resultNK.WindPtsR1,
    A_resultNK.GateR1,
    A_resultNK.TotRun1,
    A_resultNK.PosCC,
    A_resultNK.TimeTot,
    A_resultNK.TimeTotInt,
    A_resultNK.PointsJump
FROM A_resultNK
INNER JOIN A_raceNK ON A_resultNK.RaceID = A_raceNK.RaceID
WHERE A_resultNK.CompetitorID = $1
  AND (A_raceNK.SeasonCode = ANY($2) OR $2 IS NULL)
  AND (A_raceNK.DisciplineCode = ANY($3) OR $3 IS NULL)
  AND (A_raceNK.CatCode = ANY($4) OR $4 IS NULL)
ORDER BY A_raceNK.RaceDate;


-- name: GetRaceResultsNKByRaceID :many
SELECT 
    RecID,
    RaceID,
    Position,
    SpeedR1,
    DistR1,
    JudptsR1,
    PosR1,
    GateR1,
    TotRun1,
    Pointsjump,
    Poscc,
    Timetot,
    Timetotint
FROM A_resultnk
WHERE (RaceID = ANY($1))
ORDER BY RaceID;


-- name: GetCompetitorIDByFiscodeJP :one
SELECT CompetitorID
FROM A_competitor
WHERE Fiscode = $1 AND SectorCode = 'JP';


-- name: GetAthleteResultsJP :many
SELECT 
    A_resultJP.RaceID,
    A_resultJP.Position,
    A_raceJP.RaceDate,
    A_raceJP.SeasonCode,
    A_raceJP.DisciplineCode,
    A_raceJP.CatCode,
    A_raceJP.Place,
    A_resultJP.PosR1,
    A_resultJP.SpeedR1,
    A_resultJP.DistR1,
    A_resultJP.JudPtsR1,
    A_resultJP.WindR1,
    A_resultJP.WindPtsR1,
    A_resultJP.GateR1,
    A_resultJP.PosR2,
    A_resultJP.SpeedR2,
    A_resultJP.DistR2,
    A_resultJP.JudPtsR2,
    A_resultJP.WindR2,
    A_resultJP.WindPtsR2,
    A_resultJP.GateR2
FROM A_resultJP
INNER JOIN A_raceJP ON A_resultJP.RaceID = A_raceJP.RaceID
WHERE A_resultJP.CompetitorID = $1
  AND (A_raceJP.SeasonCode = ANY($2) OR $2 IS NULL)
  AND (A_raceJP.DisciplineCode = ANY($3) OR $3 IS NULL)
  AND (A_raceJP.CatCode = ANY($4) OR $4 IS NULL)
ORDER BY A_raceJP.RaceDate;



-- name: GetRaceResultsJPByRaceID :many
SELECT 
    RecID,
    RaceID,
    Position,
    SpeedR1,
    DistR1,
    JudptsR1,
    PosR1,
    GateR1,
    SpeedR2,
    DistR2,
    JudptsR2,
    PosR2,
    GateR2
FROM A_resultjp
WHERE (RaceID = ANY($1))
ORDER BY RaceID;


-- name: GetCompetitorIDsByGenderAndNationJP :many
SELECT CompetitorID
FROM A_competitor
WHERE Gender = $1 AND NationCode = $2 AND SectorCode = 'JP';


-- name: GetResultsByCompetitorsJP :many
SELECT 
    A_resultJP.RaceID,
    A_resultJP.Position,
    A_raceJP.RaceDate,
    A_raceJP.SeasonCode,
    A_raceJP.DisciplineCode,
    A_raceJP.CatCode,
    A_raceJP.Place,
    A_resultJP.PosR1,
    A_resultJP.SpeedR1,
    A_resultJP.DistR1,
    A_resultJP.JudPtsR1,
    A_resultJP.WindR1,
    A_resultJP.WindPtsR1,
    A_resultJP.GateR1,
    A_resultJP.PosR2,
    A_resultJP.SpeedR2,
    A_resultJP.DistR2,
    A_resultJP.JudPtsR2,
    A_resultJP.WindR2,
    A_resultJP.WindPtsR2,
    A_resultJP.GateR2
FROM A_resultJP
INNER JOIN A_raceJP ON A_resultJP.RaceID = A_raceJP.RaceID
WHERE A_resultJP.CompetitorID = ANY($1)
  AND (A_raceJP.SeasonCode = ANY($2) OR $2 IS NULL)
  AND (A_raceJP.DisciplineCode = ANY($3) OR $3 IS NULL)
  AND (A_raceJP.CatCode = ANY($4) OR $4 IS NULL)
ORDER BY A_raceJP.RaceDate;

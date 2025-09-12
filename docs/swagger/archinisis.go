package swagger

type UserDataArchinisisResponse struct {
	Users []string `json:"sportti_ids" example:"12345,67890,54321"`
}

type RaceReportSessionsResponse struct {
	RaceReport []int32 `json:"race_report" example:"101,102,103"`
}

type ArchRaceReportUpsertRequest struct {
	SporttiID  string `json:"sportti_id" example:"27578816"`
	SessionID  int32  `json:"session_id" example:"1842"`
	RaceReport string `json:"race_report" example:"<!DOCTYPE html><html><head><title>Race</title></head><body><h1>Report</h1></body></html>"`
}

type ArchDataUpsertRequest struct {
	NationalID   string                  `json:"national_id" example:"27353728"`
	FirstName    *string                 `json:"first_name,omitempty" example:"Jane"`
	LastName     *string                 `json:"last_name,omitempty" example:"Doe"`
	Initials     *string                 `json:"initials,omitempty" example:"JD"`
	DateOfBirth  *string                 `json:"date_of_birth,omitempty" example:"1995-07-21"`
	Height       *float64                `json:"height,omitempty" example:"1.72"`
	Weight       *float64                `json:"weight,omitempty" example:"65.4"`
	Measurements []ArchMeasurementUpsert `json:"measurements,omitempty"`
}

type ArchMeasurementUpsert struct {
	MeasurementGroupID int32   `json:"measurement_group_id" example:"2024010101"`
	MeasurementID      int32   `json:"measurement_id" example:"42"`
	Discipline         *string `json:"discipline,omitempty" example:"Sprint 100m"`
	SessionName        *string `json:"session_name,omitempty" example:"Morning session"`
	Place              *string `json:"place,omitempty" example:"Helsinki"`
	RaceID             *int32  `json:"race_id,omitempty" example:"123"`
	StartTime          *string `json:"start_time,omitempty" example:"2024-01-01T09:00:00Z"`
	StopTime           *string `json:"stop_time,omitempty" example:"2024-01-01T09:15:00Z"`
	NbSegments         *int32  `json:"nb_segments,omitempty" example:"5"`
	Comment            *string `json:"comment,omitempty" example:"Felt strong today."`
}

type ArchDataResponse struct {
	NationalID   string                    `json:"national_id" example:"27353728"`
	FirstName    *string                   `json:"first_name,omitempty" example:"Jane"`
	LastName     *string                   `json:"last_name,omitempty" example:"Doe"`
	Initials     *string                   `json:"initials,omitempty" example:"JD"`
	DateOfBirth  *string                   `json:"date_of_birth,omitempty" example:"1995-07-21"`
	Height       *float64                  `json:"height,omitempty" example:"1.72"`
	Weight       *float64                  `json:"weight,omitempty" example:"65.4"`
	Measurements []ArchMeasurementResponse `json:"measurements"`
}

type ArchMeasurementResponse struct {
	MeasurementGroupID int32   `json:"measurement_group_id" example:"2024010101"`
	MeasurementID      *int32  `json:"measurement_id,omitempty" example:"42"`
	Discipline         *string `json:"discipline,omitempty" example:"Sprint 100m"`
	SessionName        *string `json:"session_name,omitempty" example:"Morning session"`
	Place              *string `json:"place,omitempty" example:"Helsinki"`
	RaceID             *int32  `json:"race_id,omitempty" example:"123"`
	StartTime          *string `json:"start_time,omitempty" example:"2024-01-01T09:00:00Z"`
	StopTime           *string `json:"stop_time,omitempty" example:"2024-01-01T09:15:00Z"`
	NbSegments         *int32  `json:"nb_segments,omitempty" example:"5"`
	Comment            *string `json:"comment,omitempty" example:"Felt strong today."`
}

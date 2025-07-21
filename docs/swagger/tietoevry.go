package swagger

// TietoevryUserUpsertInput represents the POST body schema for upserting a user
type TietoevryUserUpsertInput struct {
	ID                        string   `json:"id" example:"7cffe6e0-3f28-43b6-b511-d836d3a9f7b5"`
	SporttiID                 int32    `json:"sportti_id" example:"12345"`
	ProfileGender             *string  `json:"profile_gender" example:"male"`
	ProfileBirthdate          *string  `json:"profile_birthdate" example:"1990-08-15"` // ISO 8601
	ProfileWeight             *float64 `json:"profile_weight" example:"72.5"`
	ProfileHeight             *float64 `json:"profile_height" example:"178.2"`
	ProfileRestingHeartRate   *int32   `json:"profile_resting_heart_rate" example:"60"`
	ProfileMaximumHeartRate   *int32   `json:"profile_maximum_heart_rate" example:"190"`
	ProfileAerobicThreshold   *int32   `json:"profile_aerobic_threshold" example:"140"`
	ProfileAnaerobicThreshold *int32   `json:"profile_anaerobic_threshold" example:"165"`
	ProfileVo2max             *int32   `json:"profile_vo2max" example:"50"`
}

type HRZone struct {
	ExerciseID    string `json:"exercise_id" example:"2d4f6aee-b62c-408e-85e1-07bd78f383a7"`
	ZoneIndex     int32  `json:"zone_index" example:"2"`
	SecondsInZone int32  `json:"seconds_in_zone" example:"90"`
	LowerLimit    int32  `json:"lower_limit" example:"120"`
	UpperLimit    int32  `json:"upper_limit" example:"140"`
	CreatedAt     string `json:"created_at" example:"2024-07-21T08:00:00Z"`
	UpdatedAt     string `json:"updated_at" example:"2024-07-21T08:00:00Z"`
}

type Sample struct {
	ID            string    `json:"id" example:"3f28c7a1-2ea3-438c-9b35-099c4372da49"`
	UserID        string    `json:"user_id" example:"7cffe6e0-3f28-43b6-b511-d836d3a9f7b5"`
	ExerciseID    string    `json:"exercise_id" example:"2d4f6aee-b62c-408e-85e1-07bd78f383a7"`
	SampleType    string    `json:"sample_type" example:"heart_rate"`
	RecordingRate int32     `json:"recording_rate" example:"1"`
	Samples       []float64 `json:"samples"`
	Source        string    `json:"source" example:"garmin"`
}

type Section struct {
	ID          string `json:"id" example:"1a2b3c4d-0000-0000-0000-111122223333"`
	UserID      string `json:"user_id" example:"7cffe6e0-3f28-43b6-b511-d836d3a9f7b5"`
	ExerciseID  string `json:"exercise_id" example:"2d4f6aee-b62c-408e-85e1-07bd78f383a7"`
	CreatedAt   string `json:"created_at" example:"2024-07-21T08:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2024-07-21T08:05:00Z"`
	StartTime   string `json:"start_time" example:"2024-07-21T08:10:00Z"`
	EndTime     string `json:"end_time" example:"2024-07-21T08:30:00Z"`
	SectionType string `json:"section_type" example:"warmup"`
	Name        string `json:"name" example:"Warm-up"`
	Comment     string `json:"comment" example:"Felt great"`
	Source      string `json:"source" example:"polar"`
	RawID       string `json:"raw_id" example:"abc-123"`
	RawData     string `json:"raw_data" example:"{\"extra\":\"data\"}"`
}

type TietoevryExerciseUpsertInput struct {
	ID                string   `json:"id" example:"2d4f6aee-b62c-408e-85e1-07bd78f383a7"`
	CreatedAt         string   `json:"created_at" example:"2024-07-21T08:00:00Z"`
	UpdatedAt         string   `json:"updated_at" example:"2024-07-21T08:00:00Z"`
	UserID            string   `json:"user_id" example:"7cffe6e0-3f28-43b6-b511-d836d3a9f7b5"`
	StartTime         string   `json:"start_time" example:"2024-07-21T08:10:00Z"`
	Duration          int64    `json:"duration" example:"3600"`
	Comment           *string  `json:"comment" example:"Morning run"`
	SportType         *string  `json:"sport_type" example:"running"`
	DetailedSportType *string  `json:"detailed_sport_type" example:"trail"`
	Distance          *float64 `json:"distance" example:"8.5"`
	AvgHeartRate      *float64 `json:"avg_heart_rate" example:"145.2"`
	MaxHeartRate      *float64 `json:"max_heart_rate" example:"165.3"`
	Trimp             *float64 `json:"trimp" example:"100.5"`
	SprintCount       *int32   `json:"sprint_count" example:"2"`
	AvgSpeed          *float64 `json:"avg_speed" example:"3.2"`
	MaxSpeed          *float64 `json:"max_speed" example:"4.8"`
	Source            string   `json:"source" example:"garmin"`
	Status            *string  `json:"status" example:"completed"`
	Calories          *int32   `json:"calories" example:"450"`
	TrainingLoad      *int32   `json:"training_load" example:"55"`
	RawID             *string  `json:"raw_id" example:"xyz789"`
	Feeling           *int32   `json:"feeling" example:"4"`
	Recovery          *int32   `json:"recovery" example:"3"`
	RPE               *int32   `json:"rpe" example:"7"`
	RawData           *string  `json:"raw_data" example:"{\"power\":300}"`

	HRZones  []HRZone  `json:"hr_zones"`
	Samples  []Sample  `json:"samples"`
	Sections []Section `json:"sections"`
}

package swagger

type KamkAddInjuryRequest struct {
	UserID      string  `json:"user_id" example:"27353728"`
	InjuryType  int32   `json:"injury_type" example:"1"`
	Severity    *int32  `json:"severity,omitempty" example:"3"`
	PainLevel   *int32  `json:"pain_level,omitempty" example:"7"`
	Description *string `json:"description,omitempty" example:"Left ankle sprain during training"`
	InjuryID    int32   `json:"injury_id" example:"2"`
	Meta        *string `json:"meta,omitempty" example:"phase=preseason"`
}

type KamkMarkRecoveredRequest struct {
	UserID   string `json:"user_id" example:"27353728"`
	InjuryID int32  `json:"injury_id" example:"2"`
}

type KamkInjuryItem struct {
	CompetitorID int32   `json:"competitor_id" example:"27353728"`
	InjuryType   int32   `json:"injury_type" example:"1"`
	Severity     *int32  `json:"severity,omitempty" example:"3"`
	PainLevel    *int32  `json:"pain_level,omitempty" example:"7"`
	Description  *string `json:"description,omitempty" example:"Left ankle sprain during training"`
	DateStart    string  `json:"date_start" example:"2025-01-10T09:30:00Z"`
	Status       int32   `json:"status" example:"0"`
	DateEnd      *string `json:"date_end,omitempty" example:"2025-01-17T12:00:00Z"`
	InjuryID     *int32  `json:"injury_id,omitempty" example:"2"`
	Meta         *string `json:"meta,omitempty" example:"phase=preseason"`
}

type KamkInjuriesListResponse struct {
	Injuries []KamkInjuryItem `json:"injuries"`
}

type KamkMaxInjuryIDResponse struct {
	ID int32 `json:"id" example:"3"`
}

package swagger

type UserDataKlabResponse struct {
	Users []string `json:"sportti_ids" example:"12345,67890,54321"`
}

// KlabCustomerResponse mirrors public.customer (51 columns)
type KlabCustomerResponse struct {
	Idcustomer         int32    `json:"idcustomer" example:"5233"`
	Firstname          string   `json:"firstname" example:"name"`
	Lastname           string   `json:"lastname" example:"lastname"`
	Idgroups           *int32   `json:"idgroups,omitempty" example:"2340"`
	Dob                *string  `json:"dob,omitempty" example:"1999-05-19"` // date (YYYY-MM-DD)
	Sex                *int32   `json:"sex,omitempty" example:"0"`
	DobYear            *int32   `json:"dob_year,omitempty" example:"1999"`
	DobMonth           *int32   `json:"dob_month,omitempty" example:"5"`
	DobDay             *int32   `json:"dob_day,omitempty" example:"9"`
	PidNumber          *string  `json:"pid_number,omitempty"`
	Company            *string  `json:"company,omitempty"`
	Occupation         *string  `json:"occupation,omitempty"`
	Education          *string  `json:"education,omitempty"`
	Address            *string  `json:"address,omitempty"`
	PhoneHome          *string  `json:"phone_home,omitempty"`
	PhoneWork          *string  `json:"phone_work,omitempty"`
	PhoneMobile        *string  `json:"phone_mobile,omitempty"`
	Faxno              *string  `json:"faxno,omitempty"`
	Email              *string  `json:"email,omitempty"`
	Username           *string  `json:"username,omitempty"`
	Password           *string  `json:"password,omitempty"`
	Readonly           *int32   `json:"readonly,omitempty"`
	Warnings           *int32   `json:"warnings,omitempty"`
	AllowToSave        *int32   `json:"allow_to_save,omitempty"`
	AllowToCloud       *int32   `json:"allow_to_cloud,omitempty"`
	Flag2              *int32   `json:"flag2,omitempty"`
	Idsport            *int32   `json:"idsport,omitempty"`
	Medication         *string  `json:"medication,omitempty"`
	Addinfo            *string  `json:"addinfo,omitempty"`
	TeamName           *string  `json:"team_name,omitempty"`
	Add1               *int32   `json:"add1,omitempty"`
	Athlete            *int32   `json:"athlete,omitempty"`
	Add10              *string  `json:"add10,omitempty"`
	Add20              *string  `json:"add20,omitempty"`
	Updatemode         *int32   `json:"updatemode,omitempty"`
	WeightKg           *float64 `json:"weight_kg,omitempty"`
	HeightCm           *float64 `json:"height_cm,omitempty"`
	DateModified       *float64 `json:"date_modified,omitempty"` // double precision
	RecomTestlevel     *int32   `json:"recom_testlevel,omitempty"`
	CreatedBy          *int64   `json:"created_by,omitempty"`
	ModBy              *int64   `json:"mod_by,omitempty"`
	ModDate            *string  `json:"mod_date,omitempty"`     // timestamp (RFC3339)
	Deleted            *int32   `json:"deleted,omitempty"`      // smallint -> int32 in JSON
	CreatedDate        *string  `json:"created_date,omitempty"` // timestamp (RFC3339)
	Modded             *int32   `json:"modded,omitempty"`       // smallint -> int32 in JSON
	AllowAnonymousData *bool    `json:"allow_anonymous_data,omitempty"`
	Locked             *int32   `json:"locked,omitempty"` // smallint -> int32
	AllowToSprintai    *int32   `json:"allow_to_sprintai,omitempty"`
	TosprintaiFrom     *string  `json:"tosprintai_from,omitempty"` // date (YYYY-MM-DD)
	StatSent           *string  `json:"stat_sent,omitempty"`       // date (YYYY-MM-DD)
	SporttiID          *string  `json:"sportti_id,omitempty"`      // varchar(64)
}

type UserKlabResponse struct {
	Customer []KlabCustomerResponse `json:"customer"`
}

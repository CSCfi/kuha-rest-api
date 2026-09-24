package archapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DeRuina/KUHA-REST-API/internal/auth/authn"
	archsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/archinisis"
	"github.com/DeRuina/KUHA-REST-API/internal/logger"
	"github.com/DeRuina/KUHA-REST-API/internal/store/archinisis"
	"go.uber.org/zap"
)

func init() {
	logger.Logger = zap.NewNop().Sugar()
}

// fakeData records calls made by the handlers.
type fakeData struct {
	archinisis.Data // unimplemented methods panic if called

	upsertDataPayload *archinisis.ArchDataPayload
	upsertDataErr     error

	getDataResp *archinisis.ArchDataResponse
	getDataErr  error

	raceReport struct {
		sporttiID  string
		sessionID  int32
		raceReport string
	}
}

func (f *fakeData) UpsertData(_ context.Context, p archinisis.ArchDataPayload) error {
	f.upsertDataPayload = &p
	return f.upsertDataErr
}

func (f *fakeData) GetDataBySporttiID(_ context.Context, _ string) (*archinisis.ArchDataResponse, error) {
	return f.getDataResp, f.getDataErr
}

func (f *fakeData) UpsertRaceReport(_ context.Context, sporttiID string, sessionID int32, raceReport string) error {
	f.raceReport.sporttiID = sporttiID
	f.raceReport.sessionID = sessionID
	f.raceReport.raceReport = raceReport
	return nil
}

func archRequest(method, target, body string) *http.Request {
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, target, nil)
	} else {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	}
	ctx := authn.WithClientMetadata(r.Context(), "test-client", []string{"archinisis"})
	return r.WithContext(ctx)
}

// Payload copied from the Archinisis endpoint listing (POST /archinisis/data).
const endpointListPayload = `{"date_of_birth": "1994-01-01", "first_name": "Petra", "height": 175.0, "initials": "PeTo", "last_name": "Torvinen", "measurements": [ {"comment": "", "discipline": "Cross-Country Skating", "measurement_group_id": 1642, "measurement_id": 6245, "nb_segments": 1, "place": "Vuokatti", "race_id": 330, "session_name": "RH SM Sprint N 2023", "start_time": "2023-09-16 05:55:24 UTC", "stop_time": "2023-09-16 06:59:42 UTC" }, {"comment": "", "discipline": "Ski Jumping", "measurement_group_id": 877, "measurement_id": 2760, "nb_segments": 8, "place": "Vuokatti", "race_id": null, "session_name": "Naos testi", "start_time": "2022-10-12 12:15:49 UTC", "stop_time": "2022-10-12 14:02:30 UTC"} ], "national_id": "27578816", "weight": 72.0}`

func TestPostArchData_EndpointListPayload(t *testing.T) {
	store := &fakeData{}
	h := NewDataHandler(store, nil)

	rec := httptest.NewRecorder()
	h.PostArchData(rec, archRequest(http.MethodPost, "/v1/archinisis/data", endpointListPayload))

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if store.upsertDataPayload == nil {
		t.Fatal("UpsertData was not called")
	}
	p := store.upsertDataPayload

	if p.Athlete.NationalID != "27578816" {
		t.Errorf("athlete national_id = %q", p.Athlete.NationalID)
	}
	if p.Athlete.FirstName.String != "Petra" || p.Athlete.LastName.String != "Torvinen" || p.Athlete.Initials.String != "PeTo" {
		t.Errorf("athlete name fields = %+v", p.Athlete)
	}
	if !p.Athlete.DateOfBirth.Valid || p.Athlete.DateOfBirth.Time.Format("2006-01-02") != "1994-01-01" {
		t.Errorf("date_of_birth = %+v", p.Athlete.DateOfBirth)
	}
	if p.Athlete.Height.String != "175" && p.Athlete.Height.String != "175.0" {
		t.Errorf("height = %+v", p.Athlete.Height)
	}

	if len(p.Measurements) != 2 {
		t.Fatalf("measurements = %d, want 2", len(p.Measurements))
	}
	m0, m1 := p.Measurements[0], p.Measurements[1]

	if m0.MeasurementID != 6245 || !m0.MeasurementGroupID.Valid || m0.MeasurementGroupID.Int32 != 1642 {
		t.Errorf("m0 ids = id %d group %+v", m0.MeasurementID, m0.MeasurementGroupID)
	}
	if m0.NationalID.String != "27578816" {
		t.Errorf("m0 national_id = %+v", m0.NationalID)
	}
	if !m0.RaceID.Valid || m0.RaceID.Int32 != 330 {
		t.Errorf("m0 race_id = %+v", m0.RaceID)
	}
	if !m0.StartTime.Valid || m0.StartTime.Time.UTC().Format("2006-01-02 15:04:05") != "2023-09-16 05:55:24" {
		t.Errorf("m0 start_time = %+v", m0.StartTime)
	}
	if !m0.StopTime.Valid || m0.StopTime.Time.UTC().Format("2006-01-02 15:04:05") != "2023-09-16 06:59:42" {
		t.Errorf("m0 stop_time = %+v", m0.StopTime)
	}

	if m1.MeasurementID != 2760 || m1.MeasurementGroupID.Int32 != 877 {
		t.Errorf("m1 ids = id %d group %+v", m1.MeasurementID, m1.MeasurementGroupID)
	}
	if m1.RaceID.Valid {
		t.Errorf("m1 race_id should be NULL, got %+v", m1.RaceID)
	}
	if !m1.NbSegments.Valid || m1.NbSegments.Int32 != 8 {
		t.Errorf("m1 nb_segments = %+v", m1.NbSegments)
	}
}

func TestPostArchData_AthleteOnly(t *testing.T) {
	store := &fakeData{}
	h := NewDataHandler(store, nil)

	rec := httptest.NewRecorder()
	h.PostArchData(rec, archRequest(http.MethodPost, "/v1/archinisis/data", `{"national_id":"27578816","first_name":"Petra"}`))

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if store.upsertDataPayload == nil || len(store.upsertDataPayload.Measurements) != 0 {
		t.Fatalf("expected UpsertData with no measurements, got %+v", store.upsertDataPayload)
	}
}

func TestPostArchData_Validation(t *testing.T) {
	cases := map[string]string{
		"non-numeric national_id":      `{"national_id":"abc"}`,
		"missing national_id":          `{"first_name":"Petra"}`,
		"measurement without group id": `{"national_id":"27578816","measurements":[{"measurement_id":1}]}`,
		"measurement without id":       `{"national_id":"27578816","measurements":[{"measurement_group_id":1}]}`,
		"bad start_time":               `{"national_id":"27578816","measurements":[{"measurement_group_id":1,"measurement_id":1,"start_time":"yesterday"}]}`,
		"unknown field":                `{"national_id":"27578816","foo":1}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			store := &fakeData{}
			h := NewDataHandler(store, nil)
			rec := httptest.NewRecorder()
			h.PostArchData(rec, archRequest(http.MethodPost, "/v1/archinisis/data", body))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
			}
			if store.upsertDataPayload != nil {
				t.Fatal("UpsertData must not be called on invalid input")
			}
		})
	}
}

func TestPostArchData_Forbidden(t *testing.T) {
	h := NewDataHandler(&fakeData{}, nil)
	r := httptest.NewRequest(http.MethodPost, "/v1/archinisis/data", strings.NewReader(endpointListPayload))
	r = r.WithContext(authn.WithClientMetadata(r.Context(), "reader", []string{"archinisis_read"}))
	rec := httptest.NewRecorder()
	h.PostArchData(rec, r)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}

func TestGetArchData(t *testing.T) {
	gid := int32(1642)
	store := &fakeData{
		getDataResp: &archinisis.ArchDataResponse{
			ArchAthleteResponse: archinisis.ArchAthleteResponse{NationalID: "27578816"},
			Measurements: []archinisis.ArchMeasurementResponse{
				{MeasurementGroupID: &gid, MeasurementID: 6245},
			},
		},
	}
	h := NewDataHandler(store, nil)

	rec := httptest.NewRecorder()
	h.GetArchData(rec, archRequest(http.MethodGet, "/v1/archinisis/data?id=27578816", ""))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var got struct {
		NationalID   string `json:"national_id"`
		Measurements []struct {
			MeasurementGroupID int32 `json:"measurement_group_id"`
			MeasurementID      int32 `json:"measurement_id"`
		} `json:"measurements"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.NationalID != "27578816" || len(got.Measurements) != 1 ||
		got.Measurements[0].MeasurementGroupID != 1642 || got.Measurements[0].MeasurementID != 6245 {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestGetArchData_NotFound(t *testing.T) {
	h := NewDataHandler(&fakeData{getDataErr: sql.ErrNoRows}, nil)
	rec := httptest.NewRecorder()
	h.GetArchData(rec, archRequest(http.MethodGet, "/v1/archinisis/data?id=1", ""))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestPostRaceReport(t *testing.T) {
	store := &fakeData{}
	h := NewDataHandler(store, nil)
	rec := httptest.NewRecorder()
	h.PostRaceReport(rec, archRequest(http.MethodPost, "/v1/archinisis/race-report",
		`{"sportti_id": "27578816", "session_id": 1642, "race_report": "<html>...</html>"}`))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if store.raceReport.sporttiID != "27578816" || store.raceReport.sessionID != 1642 || store.raceReport.raceReport != "<html>...</html>" {
		t.Errorf("UpsertRaceReport args = %+v", store.raceReport)
	}
}

func TestMapMeasurementToParams_GroupIDIsNullableInt(t *testing.T) {
	p, err := mapMeasurementToParams(ArchMeasurementInput{MeasurementGroupID: 877, MeasurementID: 2760}, "1")
	if err != nil {
		t.Fatal(err)
	}
	if p.MeasurementID != 2760 {
		t.Errorf("measurement_id = %d", p.MeasurementID)
	}
	if !p.MeasurementGroupID.Valid || p.MeasurementGroupID.Int32 != 877 {
		t.Errorf("measurement_group_id = %+v", p.MeasurementGroupID)
	}
	var _ archsqlc.UpsertMeasurementParams = p
}

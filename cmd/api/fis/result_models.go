package fisapi

import (
	"time"

	fissqlc "github.com/DeRuina/KUHA-REST-API/internal/db/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/store/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type AthleteResultsCCQuery struct {
	SeasonCode     []string `form:"seasoncode"`
	DisciplineCode []string `form:"disciplinecode"`
	CatCode        []string `form:"catcode"`
}

type InsertResultCCInput struct {
	Recid          int32   `json:"recid" validate:"required"`
	Raceid         *int32  `json:"raceid"`
	Competitorid   *int32  `json:"competitorid"`
	Status         *string `json:"status"`
	Reason         *string `json:"reason"`
	Position       *string `json:"position"`
	Pf             *int32  `json:"pf"`
	Status2        *string `json:"status2"`
	Bib            *string `json:"bib"`
	Bibcolor       *string `json:"bibcolor"`
	Fiscode        *int32  `json:"fiscode"`
	Competitorname *string `json:"competitorname"`
	Nationcode     *string `json:"nationcode"`
	Stage          *string `json:"stage"`
	Level          *string `json:"level"`
	Heat           *string `json:"heat"`
	Timer1         *string `json:"timer1"`
	Timer2         *string `json:"timer2"`
	Timer3         *string `json:"timer3"`
	Timetot        *string `json:"timetot"`
	Valid          *string `json:"valid"`
	Racepoints     *string `json:"racepoints"`
	Cuppoints      *string `json:"cuppoints"`
	Bonustime      *string `json:"bonustime"`
	Bonuscuppoints *string `json:"bonuscuppoints"`
	Version        *string `json:"version"`
	Rg1            *string `json:"rg1"`
	Rg2            *string `json:"rg2"`
	Lastupdate     *string `json:"lastupdate"` // RFC3339
}

type UpdateResultCCInput = InsertResultCCInput

func mapInsertResultCCInput(in InsertResultCCInput) (fis.InsertResultCCClean, error) {
	var lup *time.Time
	var err error

	if in.Lastupdate != nil {
		lup, err = utils.ParseTimestampPtr(in.Lastupdate)
		if err != nil {
			return fis.InsertResultCCClean{}, err
		}
	}

	return fis.InsertResultCCClean{
		Recid:          in.Recid,
		Raceid:         in.Raceid,
		Competitorid:   in.Competitorid,
		Status:         in.Status,
		Reason:         in.Reason,
		Position:       in.Position,
		Pf:             in.Pf,
		Status2:        in.Status2,
		Bib:            in.Bib,
		Bibcolor:       in.Bibcolor,
		Fiscode:        in.Fiscode,
		Competitorname: in.Competitorname,
		Nationcode:     in.Nationcode,
		Stage:          in.Stage,
		Level:          in.Level,
		Heat:           in.Heat,
		Timer1:         in.Timer1,
		Timer2:         in.Timer2,
		Timer3:         in.Timer3,
		Timetot:        in.Timetot,
		Valid:          in.Valid,
		Racepoints:     in.Racepoints,
		Cuppoints:      in.Cuppoints,
		Bonustime:      in.Bonustime,
		Bonuscuppoints: in.Bonuscuppoints,
		Version:        in.Version,
		Rg1:            in.Rg1,
		Rg2:            in.Rg2,
		Lastupdate:     lup,
	}, nil
}

func mapUpdateResultCCInput(in UpdateResultCCInput) (fis.UpdateResultCCClean, error) {
	clean, err := mapInsertResultCCInput(InsertResultCCInput(in))
	return fis.UpdateResultCCClean(clean), err
}

type FISResultCCFullResponse struct {
	Recid          int32   `json:"recid"`
	Raceid         *int32  `json:"raceid"`
	Competitorid   *int32  `json:"competitorid"`
	Status         *string `json:"status"`
	Reason         *string `json:"reason"`
	Position       *string `json:"position"`
	Pf             *int32  `json:"pf"`
	Status2        *string `json:"status2"`
	Bib            *string `json:"bib"`
	Bibcolor       *string `json:"bibcolor"`
	Fiscode        *int32  `json:"fiscode"`
	Competitorname *string `json:"competitorname"`
	Nationcode     *string `json:"nationcode"`
	Stage          *string `json:"stage"`
	Level          *string `json:"level"`
	Heat           *string `json:"heat"`
	Timer1         *string `json:"timer1"`
	Timer2         *string `json:"timer2"`
	Timer3         *string `json:"timer3"`
	Timetot        *string `json:"timetot"`
	Valid          *string `json:"valid"`
	Racepoints     *string `json:"racepoints"`
	Cuppoints      *string `json:"cuppoints"`
	Bonustime      *string `json:"bonustime"`
	Bonuscuppoints *string `json:"bonuscuppoints"`
	Version        *string `json:"version"`
	Rg1            *string `json:"rg1"`
	Rg2            *string `json:"rg2"`
	Lastupdate     *string `json:"lastupdate"`
}

func FISResultCCFullFromSqlc(row fissqlc.AResultcc) FISResultCCFullResponse {
	var lastUpdateStr *string
	if row.Lastupdate.Valid {
		lastUpdateStr = utils.FormatTimestampPtr(row.Lastupdate)
	}

	return FISResultCCFullResponse{
		Recid:          row.Recid,
		Raceid:         utils.Int32PtrOrNil(row.Raceid),
		Competitorid:   utils.Int32PtrOrNil(row.Competitorid),
		Status:         utils.StringPtrOrNil(row.Status),
		Reason:         utils.StringPtrOrNil(row.Reason),
		Position:       utils.StringPtrOrNil(row.Position),
		Pf:             utils.Int32PtrOrNil(row.Pf),
		Status2:        utils.StringPtrOrNil(row.Status2),
		Bib:            utils.StringPtrOrNil(row.Bib),
		Bibcolor:       utils.StringPtrOrNil(row.Bibcolor),
		Fiscode:        utils.Int32PtrOrNil(row.Fiscode),
		Competitorname: utils.StringPtrOrNil(row.Competitorname),
		Nationcode:     utils.StringPtrOrNil(row.Nationcode),
		Stage:          utils.StringPtrOrNil(row.Stage),
		Level:          utils.StringPtrOrNil(row.Level),
		Heat:           utils.StringPtrOrNil(row.Heat),
		Timer1:         utils.StringPtrOrNil(row.Timer1),
		Timer2:         utils.StringPtrOrNil(row.Timer2),
		Timer3:         utils.StringPtrOrNil(row.Timer3),
		Timetot:        utils.StringPtrOrNil(row.Timetot),
		Valid:          utils.StringPtrOrNil(row.Valid),
		Racepoints:     utils.StringPtrOrNil(row.Racepoints),
		Cuppoints:      utils.StringPtrOrNil(row.Cuppoints),
		Bonustime:      utils.StringPtrOrNil(row.Bonustime),
		Bonuscuppoints: utils.StringPtrOrNil(row.Bonuscuppoints),
		Version:        utils.StringPtrOrNil(row.Version),
		Rg1:            utils.StringPtrOrNil(row.Rg1),
		Rg2:            utils.StringPtrOrNil(row.Rg2),
		Lastupdate:     lastUpdateStr,
	}
}

type FISAthleteResultCCRow struct {
	Recid          int32   `json:"recid"`
	Raceid         *int32  `json:"raceid"`
	Position       *string `json:"position"`
	Timetot        *string `json:"timetot"`
	Competitorid   *int32  `json:"competitorid"`
	Racedate       *string `json:"racedate"`
	Seasoncode     *int32  `json:"seasoncode"`
	Disciplinecode *string `json:"disciplinecode"`
	Catcode        *string `json:"catcode"`
	Place          *string `json:"place"`
}

func FISAthleteResultCCFromSqlc(row fissqlc.GetAthleteResultsCCRow) FISAthleteResultCCRow {
	var raceDateStr *string
	if row.Racedate.Valid {
		raceDateStr = utils.FormatDatePtr(row.Racedate)
	}
	return FISAthleteResultCCRow{
		Recid:          row.Recid,
		Raceid:         utils.Int32PtrOrNil(row.Raceid),
		Position:       utils.StringPtrOrNil(row.Position),
		Timetot:        utils.StringPtrOrNil(row.Timetot),
		Competitorid:   utils.Int32PtrOrNil(row.Competitorid),
		Racedate:       raceDateStr,
		Seasoncode:     utils.Int32PtrOrNil(row.Seasoncode),
		Disciplinecode: utils.StringPtrOrNil(row.Disciplinecode),
		Catcode:        utils.StringPtrOrNil(row.Catcode),
		Place:          utils.StringPtrOrNil(row.Place),
	}
}

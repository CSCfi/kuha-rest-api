package fis

import (
	"time"

	fissqlc "github.com/DeRuina/KUHA-REST-API/internal/db/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

// Clean structs for INSERT/UPDATE on a_resultcc

type InsertResultCCClean struct {
	Recid          int32
	Raceid         *int32
	Competitorid   *int32
	Status         *string
	Reason         *string
	Position       *string // numeric(10,5) in DB, but handled as string
	Pf             *int32
	Status2        *string
	Bib            *string // numeric(10,5) in DB, but handled as string
	Bibcolor       *string
	Fiscode        *int32
	Competitorname *string
	Nationcode     *string
	Stage          *string
	Level          *string
	Heat           *string
	Timer1         *string
	Timer2         *string
	Timer3         *string
	Timetot        *string
	Valid          *string // numeric(10,5) in DB, but handled as string
	Racepoints     *string // numeric(10,5) in DB, but handled as string
	Cuppoints      *string // numeric(10,5) in DB, but handled as string
	Bonustime      *string
	Bonuscuppoints *string
	Version        *string
	Rg1            *string
	Rg2            *string
	Lastupdate     *time.Time
}

type UpdateResultCCClean = InsertResultCCClean

func mapInsertResultCCToParams(in InsertResultCCClean) fissqlc.InsertResultCCParams {
	return fissqlc.InsertResultCCParams{
		Recid:          in.Recid,
		Raceid:         utils.NullInt32Ptr(in.Raceid),
		Competitorid:   utils.NullInt32Ptr(in.Competitorid),
		Status:         utils.NullStringPtr(in.Status),
		Reason:         utils.NullStringPtr(in.Reason),
		Position:       utils.NullStringPtr(in.Position),
		Pf:             utils.NullInt32Ptr(in.Pf),
		Status2:        utils.NullStringPtr(in.Status2),
		Bib:            utils.NullStringPtr(in.Bib),
		Bibcolor:       utils.NullStringPtr(in.Bibcolor),
		Fiscode:        utils.NullInt32Ptr(in.Fiscode),
		Competitorname: utils.NullStringPtr(in.Competitorname),
		Nationcode:     utils.NullStringPtr(in.Nationcode),
		Stage:          utils.NullStringPtr(in.Stage),
		Level:          utils.NullStringPtr(in.Level),
		Heat:           utils.NullStringPtr(in.Heat),
		Timer1:         utils.NullStringPtr(in.Timer1),
		Timer2:         utils.NullStringPtr(in.Timer2),
		Timer3:         utils.NullStringPtr(in.Timer3),
		Timetot:        utils.NullStringPtr(in.Timetot),
		Valid:          utils.NullStringPtr(in.Valid),
		Racepoints:     utils.NullStringPtr(in.Racepoints),
		Cuppoints:      utils.NullStringPtr(in.Cuppoints),
		Bonustime:      utils.NullStringPtr(in.Bonustime),
		Bonuscuppoints: utils.NullStringPtr(in.Bonuscuppoints),
		Version:        utils.NullStringPtr(in.Version),
		Rg1:            utils.NullStringPtr(in.Rg1),
		Rg2:            utils.NullStringPtr(in.Rg2),
		Lastupdate:     utils.NullTimePtr(in.Lastupdate),
	}
}

func mapUpdateResultCCToParams(in UpdateResultCCClean) fissqlc.UpdateResultCCByRecIDParams {
	p := mapInsertResultCCToParams(InsertResultCCClean(in))
	return fissqlc.UpdateResultCCByRecIDParams{
		Recid:          p.Recid,
		Raceid:         p.Raceid,
		Competitorid:   p.Competitorid,
		Status:         p.Status,
		Reason:         p.Reason,
		Position:       p.Position,
		Pf:             p.Pf,
		Status2:        p.Status2,
		Bib:            p.Bib,
		Bibcolor:       p.Bibcolor,
		Fiscode:        p.Fiscode,
		Competitorname: p.Competitorname,
		Nationcode:     p.Nationcode,
		Stage:          p.Stage,
		Level:          p.Level,
		Heat:           p.Heat,
		Timer1:         p.Timer1,
		Timer2:         p.Timer2,
		Timer3:         p.Timer3,
		Timetot:        p.Timetot,
		Valid:          p.Valid,
		Racepoints:     p.Racepoints,
		Cuppoints:      p.Cuppoints,
		Bonustime:      p.Bonustime,
		Bonuscuppoints: p.Bonuscuppoints,
		Version:        p.Version,
		Rg1:            p.Rg1,
		Rg2:            p.Rg2,
		Lastupdate:     p.Lastupdate,
	}
}

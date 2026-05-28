package archapi

import (
	"context"
	"fmt"
	"time"

	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
)

const (
	ARCHCacheTTL = 6 * time.Hour

	archSessionsPrefix     = "arch:race-report:sessions"
	archHTMLPrefix         = "arch:race-report:html"
	archAthletePrefix      = "arch:athlete"
	archMeasurementsPrefix = "arch:measurements"
)

func invalidateArchRaceReport(ctx context.Context, c *cache.Storage, sporttiID string, sessionID *int32) {
	if c == nil {
		return
	}
	pf := []string{
		fmt.Sprintf("%s:%s", archSessionsPrefix, sporttiID),
	}
	if sessionID != nil {
		pf = append(pf, fmt.Sprintf("%s:%s:%d", archHTMLPrefix, sporttiID, *sessionID))
	} else {
		pf = append(pf, fmt.Sprintf("%s:%s:", archHTMLPrefix, sporttiID))
	}
	_ = c.DeleteByPrefixes(ctx, pf...)
}

func invalidateArchAthlete(ctx context.Context, c *cache.Storage, sporttiID string) {
	if c == nil {
		return
	}
	_ = c.DeleteByPrefixes(ctx, fmt.Sprintf("%s:%s", archAthletePrefix, sporttiID))
}

func invalidateArchMeasurements(ctx context.Context, c *cache.Storage, sporttiID string) {
	if c == nil {
		return
	}
	_ = c.DeleteByPrefixes(ctx, fmt.Sprintf("%s:%s", archMeasurementsPrefix, sporttiID))
}

func invalidateArchAll(ctx context.Context, c *cache.Storage, sporttiID string) {
	if c == nil {
		return
	}
	_ = c.DeleteByPrefixes(
		ctx,
		fmt.Sprintf("%s:%s", archAthletePrefix, sporttiID),
		fmt.Sprintf("%s:%s", archMeasurementsPrefix, sporttiID),
		fmt.Sprintf("%s:%s", archSessionsPrefix, sporttiID),
		fmt.Sprintf("%s:%s:", archHTMLPrefix, sporttiID),
	)
}

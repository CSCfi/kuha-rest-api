package tietoevryapi

import (
	"context"
	"fmt"
	"time"

	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
	"github.com/google/uuid"
)

const (
	TietoevryCacheTTL = 6 * time.Hour

	tzPrefix = "tietoevry:activity-zones" // Activity Zones
	exPrefix = "tietoevry:exercises"      // Exercises
	msPrefix = "tietoevry:measurements"   // Measurements
	qnPrefix = "tietoevry:questionnaires" // Questionnaires
	syPrefix = "tietoevry:symptoms"       // Symptoms
	trPrefix = "tietoevry:test-results"   // Test Results
)

// invalidate all cached variants for these resources for a user
func invalidateTietoevry(ctx context.Context, c *cache.Storage, userID uuid.UUID, resources ...string) {
	if c == nil {
		return
	}
	if len(resources) == 0 {
		resources = []string{tzPrefix, exPrefix, msPrefix, qnPrefix, syPrefix, trPrefix}
	}
	for _, r := range resources {
		_ = c.DeleteByPattern(ctx, fmt.Sprintf("%s:%s*", r, userID.String()))
	}
}

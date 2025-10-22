package klabapi

import (
	"context"
	"fmt"
	"time"

	"github.com/DeRuina/KUHA-REST-API/internal/store/cache"
)

const (
	KLABCacheTTL = 6 * time.Hour

	klabUserPrefix = "klab:user"
	klabDataPrefix = "klab:data"
	klabMapPrefix  = "klab:cid"
)

func invalidateKlabAll(ctx context.Context, c *cache.Storage, sporttiID string) {
	if c == nil {
		return
	}
	_ = c.DeleteByPrefixes(
		ctx,
		fmt.Sprintf("%s:%s", klabUserPrefix, sporttiID),
		fmt.Sprintf("%s:%s", klabDataPrefix, sporttiID),
		fmt.Sprintf("%s:%s", klabMapPrefix, sporttiID),
	)
}

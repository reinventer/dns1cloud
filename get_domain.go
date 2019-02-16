package dns1cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GetDomain returns domain by id
func (c *DNS1Cloud) GetDomain(ctx context.Context, domainID uint64) (Domain, error) {
	cmd := command{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("dns/%d", domainID),
	}

	var domain Domain
	if err := c.send(ctx, cmd, &domain); err != nil {
		return domain, errors.Wrap(err, "could not send command get_domain")
	}
	return domain, nil
}

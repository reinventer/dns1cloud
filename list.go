package dns1cloud

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// List returns list of domains
func (c *DNS1Cloud) List(ctx context.Context) ([]Domain, error) {
	cmd := command{
		method:   http.MethodGet,
		endpoint: "dns",
	}

	var domains []Domain
	if err := c.send(ctx, cmd, &domains); err != nil {
		return nil, errors.Wrap(err, "could not send command list")
	}
	return domains, nil
}

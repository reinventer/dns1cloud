package dns1cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// DeleteRecord deletes record from domain
func (c *DNS1Cloud) DeleteRecord(
	ctx context.Context,
	domainID uint64,
	recordID uint64,
) error {
	cmd := command{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("dns/%d/%d", domainID, recordID),
	}
	if err := c.send(ctx, cmd, nil); err != nil {
		return errors.Wrap(err, "could not send command delete_record")
	}
	return nil
}

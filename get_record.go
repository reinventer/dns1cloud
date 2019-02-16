package dns1cloud

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GetRecord returns record by id
func (c *DNS1Cloud) GetRecord(ctx context.Context, recordID uint64) (Record, error) {
	cmd := command{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("dns/record/%d", recordID),
	}

	var record Record
	if err := c.send(ctx, cmd, &record); err != nil {
		return record, errors.Wrap(err, "could not send command get_record")
	}
	return record, nil
}

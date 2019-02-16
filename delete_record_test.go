package dns1cloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNS1Cloud_DeleteRecord(t *testing.T) {
	testCases := []struct {
		name           string
		responseStatus int
		expErrString   string
	}{
		{
			name:           "success",
			responseStatus: http.StatusOK,
		},
		{
			name:           "bad response",
			responseStatus: http.StatusInternalServerError,
			expErrString:   `could not send command delete_record: bad response, status: 500, body: ''`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/dns/123/124", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)
				w.WriteHeader(tc.responseStatus)
			}))
			defer s.Close()

			c := New("apiKey", WithApiHost(s.URL))

			err := c.DeleteRecord(context.Background(), 123, 124)
			if len(tc.expErrString) > 0 {
				assert.EqualError(t, err, tc.expErrString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

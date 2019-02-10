package dns1cloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDNS1Cloud_List(t *testing.T) {
	testCases := []struct {
		name           string
		responseStatus int
		responseJSON   string
		expList        []Domain
		expErrString   string
	}{
		{
			name:           "success",
			responseStatus: http.StatusOK,
			responseJSON:   `[{"ID":123,"Name":"domain.com","TechName":"domain_com","State":"Active","DateCreate":"2019-02-02T21:17:11.64","IsDelegate":false,"LinkedRecords":[{"ID":124,"TypeRecord":"SRV","IP":"","HostName":"@","Priority":"20","Text":"","MnemonicName":"","ExtHostName":"","State":"Active","DateCreate":"2019-02-02T21:18:06.743","Service":"_xmpp-client.","Proto":"tcp","Weight":"0","TTL":21160,"Port":"5222","Target":"domain-xmpp.test.com.","CanonicalDescription":"_xmpp-client._tcp.domain.com. 21160 IN SRV 20 0 5222 domain-xmpp.test.com."},{"ID":125,"TypeRecord":"SRV","IP":"","HostName":"@","Priority":"20","Text":"","MnemonicName":"","ExtHostName":"","State":"Active","DateCreate":"2019-02-02T21:18:06.743","Service":"_xmpp-client.","Proto":"tcp","Weight":"0","TTL":21160,"Port":"5222","Target":"domain-xmpp.test.net.","CanonicalDescription":"_xmpp-client._tcp.domain.com. 21160 IN SRV 20 0 5222 domain-xmpp.test.net."},{"ID":126,"TypeRecord":"SRV","IP":"","HostName":"@","Priority":"20","Text":"","MnemonicName":"","ExtHostName":"","State":"Active","DateCreate":"2019-02-02T21:18:06.743","Service":"_xmpp-server.","Proto":"tcp","Weight":"0","TTL":21160,"Port":"5269","Target":"domain-xmpp.test.net.","CanonicalDescription":"_xmpp-server._tcp.domain.com. 21160 IN SRV 20 0 5269 domain-xmpp.test.net."}]}]`,
			expList: []Domain{
				{
					ID:         123,
					Name:       "domain.com",
					TechName:   "domain_com",
					State:      StateActive,
					DateCreate: DateTime{Time: time.Date(2019, 2, 2, 21, 17, 11, 640000000, time.UTC)},
					IsDelegate: false,
					LinkedRecords: []Record{
						{
							ID:                   124,
							TypeRecord:           RecordTypeSRV,
							IP:                   "",
							HostName:             "@",
							Priority:             "20",
							Text:                 "",
							MnemonicName:         "",
							ExtHostName:          "",
							Weight:               "0",
							Port:                 "5222",
							Target:               "domain-xmpp.test.com.",
							Proto:                "tcp",
							Service:              "_xmpp-client.",
							TTL:                  21160,
							State:                StateActive,
							DateCreate:           DateTime{Time: time.Date(2019, 2, 2, 21, 18, 6, 743000000, time.UTC)},
							CanonicalDescription: "_xmpp-client._tcp.domain.com. 21160 IN SRV 20 0 5222 domain-xmpp.test.com.",
						},
						{
							ID:                   125,
							TypeRecord:           RecordTypeSRV,
							IP:                   "",
							HostName:             "@",
							Priority:             "20",
							Text:                 "",
							MnemonicName:         "",
							ExtHostName:          "",
							Weight:               "0",
							Port:                 "5222",
							Target:               "domain-xmpp.test.net.",
							Proto:                "tcp",
							Service:              "_xmpp-client.",
							TTL:                  21160,
							State:                StateActive,
							DateCreate:           DateTime{Time: time.Date(2019, 2, 2, 21, 18, 6, 743000000, time.UTC)},
							CanonicalDescription: "_xmpp-client._tcp.domain.com. 21160 IN SRV 20 0 5222 domain-xmpp.test.net.",
						},
						{
							ID:                   126,
							TypeRecord:           RecordTypeSRV,
							IP:                   "",
							HostName:             "@",
							Priority:             "20",
							Text:                 "",
							MnemonicName:         "",
							ExtHostName:          "",
							Weight:               "0",
							Port:                 "5269",
							Target:               "domain-xmpp.test.net.",
							Proto:                "tcp",
							Service:              "_xmpp-server.",
							TTL:                  21160,
							State:                StateActive,
							DateCreate:           DateTime{Time: time.Date(2019, 2, 2, 21, 18, 6, 743000000, time.UTC)},
							CanonicalDescription: "_xmpp-server._tcp.domain.com. 21160 IN SRV 20 0 5269 domain-xmpp.test.net.",
						},
					},
				},
			},
		},
		{
			name:           "invalid json",
			responseStatus: http.StatusOK,
			responseJSON:   "invalid json",
			expErrString:   "could not send command list: could not unmarshal response: invalid character 'i' looking for beginning of value",
		},
		{
			name:           "incorrect json",
			responseStatus: http.StatusOK,
			responseJSON:   `{"foo": "bar"}`,
			expErrString:   "could not send command list: could not unmarshal response: json: cannot unmarshal object into Go value of type []dns1cloud.Domain",
		},
		{
			name:           "bad response",
			responseStatus: http.StatusInternalServerError,
			responseJSON:   `{"Message": "oops, error"}`,
			expErrString:   "could not send command list: bad response, status: 500, body: {\"Message\": \"oops, error\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/dns", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.responseJSON))
			}))
			defer s.Close()

			c := New("apiKey", WithApiHost(s.URL))

			list, err := c.List(context.Background())
			if len(tc.expErrString) > 0 {
				assert.EqualError(t, err, tc.expErrString)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expList, list)
		})
	}
}

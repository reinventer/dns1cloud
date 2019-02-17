package dns1cloud

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNS1Cloud_AddRecord(t *testing.T) {
	domainID := uint64(123)

	testCases := []struct {
		name           string
		reqRecord      Record
		responseStatus int
		responseJSON   string
		expPath        string
		expRequest     string
		expRecord      Record
		expErrString   string
	}{
		{
			name:           "success add record A",
			reqRecord:      Record{TypeRecord: RecordTypeA, HostName: "@", IP: "1.1.1.2", TTL: 3600},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "A", "IP": "1.1.1.2", "HostName": "@", "Priority": "", "Text": "",
				"MnemonicName": "", "ExtHostName": "","State": "Active", "TTL": 3600,
				"CanonicalDescription": "@ 3600 IN A 1.1.1.2"}`,
			expPath:    "/dns/recorda",
			expRequest: `{"DomainId": "123", "IP": "1.1.1.2", "Name": "@", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeA,
				HostName:             "@",
				IP:                   "1.1.1.2",
				State:                StateActive,
				TTL:                  3600,
				CanonicalDescription: "@ 3600 IN A 1.1.1.2",
			},
		},
		{
			name:           "success add record A with empty ttl",
			reqRecord:      Record{TypeRecord: RecordTypeA, HostName: "@", IP: "1.1.1.2", TTL: 0},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "A", "IP": "1.1.1.2", "HostName": "@", "Priority": "", "Text": "",
				"MnemonicName": "", "ExtHostName": "","State": "Active", "TTL": 0,
				"CanonicalDescription": "@ 3600 IN A 1.1.1.2"}`,
			expPath:    "/dns/recorda",
			expRequest: `{"DomainId": "123", "IP": "1.1.1.2", "Name": "@"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeA,
				HostName:             "@",
				IP:                   "1.1.1.2",
				State:                StateActive,
				TTL:                  0,
				CanonicalDescription: "@ 3600 IN A 1.1.1.2",
			},
		},
		{
			name:         "incorrect ttl for record A",
			reqRecord:    Record{TypeRecord: RecordTypeA, HostName: "@", IP: "1.1.1.2", TTL: 7},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name:         "incorrect ip for record A",
			reqRecord:    Record{TypeRecord: RecordTypeA, HostName: "@", IP: "ip", TTL: 300},
			expRecord:    Record{},
			expErrString: `IP "ip" is incorrect`,
		},
		{
			name:           "success add record AAAA",
			reqRecord:      Record{TypeRecord: RecordTypeAAAA, HostName: "@", IP: "2001:db8::68", TTL: 3600},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "AAAA", "IP": "2001:db8::68", "HostName": "@", "Priority": "", "Text": "",
				"MnemonicName": "", "ExtHostName": "","State": "Active", "TTL": 3600,
				"CanonicalDescription": "@ 3600 IN AAAA 2001:db8::68"}`,
			expPath:    "/dns/recordaaaa",
			expRequest: `{"DomainId": "123", "IP": "2001:db8::68", "Name": "@", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeAAAA,
				HostName:             "@",
				IP:                   "2001:db8::68",
				State:                StateActive,
				TTL:                  3600,
				CanonicalDescription: "@ 3600 IN AAAA 2001:db8::68",
			},
		},
		{
			name:         "incorrect ttl for record AAAA",
			reqRecord:    Record{TypeRecord: RecordTypeAAAA, HostName: "@", IP: "2001:db8::68", TTL: 7},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name:         "incorrect ip for record AAAA",
			reqRecord:    Record{TypeRecord: RecordTypeAAAA, HostName: "@", IP: "ip", TTL: 300},
			expRecord:    Record{},
			expErrString: `IP "ip" is incorrect`,
		},
		{
			name:           "success add record CNAME",
			reqRecord:      Record{TypeRecord: RecordTypeCNAME, HostName: "@", MnemonicName: "test", TTL: 3600},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "CNAME", "IP": "", "HostName": "@", "Priority": "", "Text": "",
				"MnemonicName": "test", "ExtHostName": "","State": "Active", "TTL": 3600,
				"CanonicalDescription": "test.test.ru. 3600 IN CNAME @"}`,
			expPath:    "/dns/recordcname",
			expRequest: `{"DomainId": "123", "Name": "@", "MnemonicName": "test", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeCNAME,
				HostName:             "@",
				MnemonicName:         "test",
				State:                StateActive,
				TTL:                  3600,
				CanonicalDescription: "test.test.ru. 3600 IN CNAME @",
			},
		},
		{
			name:         "incorrect ttl for record CNAME",
			reqRecord:    Record{TypeRecord: RecordTypeCNAME, HostName: "@", MnemonicName: "test", TTL: 7},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name:           "success add record MX",
			reqRecord:      Record{TypeRecord: RecordTypeMX, HostName: "mail.test.com", Priority: "10", TTL: 3600},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "MX", "IP": "", "HostName": "mail.test.com", "Priority": "10", "Text": "",
				"MnemonicName": "", "ExtHostName": "","State": "New", "TTL": 3600,
				"CanonicalDescription": "@ 3600 MX 10 mail.test.com."}`,
			expPath:    "/dns/recordmx",
			expRequest: `{"DomainId": "123", "HostName": "mail.test.com", "Priority": "10", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeMX,
				HostName:             "mail.test.com",
				MnemonicName:         "",
				State:                StateNew,
				Priority:             "10",
				TTL:                  3600,
				CanonicalDescription: "@ 3600 MX 10 mail.test.com.",
			},
		},
		{
			name:         "incorrect ttl for record MX",
			reqRecord:    Record{TypeRecord: RecordTypeMX, HostName: "@", Priority: "10", TTL: 7},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name:           "success add record NS",
			reqRecord:      Record{TypeRecord: RecordTypeNS, HostName: "ns.test.com", ExtHostName: "ext.test.com", TTL: 3600},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "NS", "IP": "", "HostName": "ns.test.com", "Priority": "", "Text": "",
				"MnemonicName": "", "ExtHostName": "ext.test.com","State": "Active", "TTL": 3600,
				"CanonicalDescription": "@ 3600 IN NS ns.test.com."}`,
			expPath:    "/dns/recordns",
			expRequest: `{"DomainId": "123", "HostName": "ns.test.com", "Name": "ext.test.com", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeNS,
				HostName:             "ns.test.com",
				ExtHostName:          "ext.test.com",
				State:                StateActive,
				Priority:             "",
				TTL:                  3600,
				CanonicalDescription: "@ 3600 IN NS ns.test.com.",
			},
		},
		{
			name:         "incorrect ttl for record NS",
			reqRecord:    Record{TypeRecord: RecordTypeNS, HostName: "ns.test.com", ExtHostName: "ext.test.com", TTL: 7},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name: "success add record SRV",
			reqRecord: Record{
				TypeRecord: RecordTypeSRV, HostName: "name.test.com", Proto: "tcp", Service: "service", Priority: "10",
				Weight: "30", Port: "4321", Target: "service.test.com", TTL: 3600,
			},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "SRV", "IP": "", "HostName": "name.test.com", "Priority": "10", "Text": "",
				"Port":"4321", "Weight": "30", "Service": "service", "MnemonicName": "", "ExtHostName": "","State": "Active", "TTL": 3600,
				"Proto": "tcp","CanonicalDescription": "_service._tcp.name.test.ru. 3600 IN SRV 1 1 4321 service.test.ru."}`,
			expPath:    "/dns/recordsrv",
			expRequest: `{"DomainId": "123", "Name": "name.test.com", "Proto":"tcp", "Service": "service", "Priority": "10", "Weight": "30", "Port": "4321", "Target": "service.test.com", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeSRV,
				HostName:             "name.test.com",
				Port:                 "4321",
				Proto:                "tcp",
				Service:              "service",
				Weight:               "30",
				State:                StateActive,
				Priority:             "10",
				TTL:                  3600,
				CanonicalDescription: "_service._tcp.name.test.ru. 3600 IN SRV 1 1 4321 service.test.ru.",
			},
		},
		{
			name: "incorrect ttl for record SRV",
			reqRecord: Record{
				TypeRecord: RecordTypeSRV, HostName: "name.test.com", Proto: "tcp", Service: "service", Priority: "10",
				Weight: "30", Port: "4321", Target: "service.test.com", TTL: 7,
			},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name:           "success add record TXT",
			reqRecord:      Record{TypeRecord: RecordTypeTXT, HostName: "text.test.com", Text: "some_text", TTL: 3600},
			responseStatus: http.StatusOK,
			responseJSON: `{"ID": 1, "TypeRecord": "TXT", "IP": "", "HostName": "text.test.com", "Text": "some_text",
				"Port":"", "Weight": "", "Service": "", "MnemonicName": "", "ExtHostName": "","State": "Active", "TTL": 3600,
				"Proto": "","CanonicalDescription": "text.test.ru. 3600 IN TXT some_text"}`,
			expPath:    "/dns/recordtxt",
			expRequest: `{"DomainId": "123", "Name": "text.test.com", "Text": "some_text", "TTL": "3600"}`,
			expRecord: Record{
				ID:                   1,
				TypeRecord:           RecordTypeTXT,
				HostName:             "text.test.com",
				Text:                 "some_text",
				State:                StateActive,
				TTL:                  3600,
				CanonicalDescription: "text.test.ru. 3600 IN TXT some_text",
			},
		},
		{
			name:         "incorrect ttl for record TXT",
			reqRecord:    Record{TypeRecord: RecordTypeTXT, HostName: "text.test.com", Text: "some_text", TTL: 7},
			expRecord:    Record{},
			expErrString: `TTL "7" is not valid`,
		},
		{
			name:         "fail: unknown record type",
			reqRecord:    Record{TypeRecord: 100},
			expRecord:    Record{},
			expErrString: `unknown record type: 100`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expPath, r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "Bearer apiKey", r.Header.Get("Authorization"))

				reqBody, err := ioutil.ReadAll(r.Body)
				defer r.Body.Close()

				assert.NoError(t, err)
				assert.JSONEq(t, tc.expRequest, string(reqBody))
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.responseJSON))
			}))
			defer s.Close()

			c := New("apiKey", WithApiHost(s.URL))

			record, err := c.AddRecord(context.Background(), domainID, tc.reqRecord)
			if len(tc.expErrString) > 0 {
				assert.EqualError(t, err, tc.expErrString)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expRecord, record)
		})
	}
}

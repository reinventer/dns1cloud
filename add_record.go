package dns1cloud

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

var validTTLs = [...]string{
	"1", "5", "30", "60", "300", "600", "900", "1800", "3600", "7200", "21160", "43200", "86400",
}

// AddRecord adds record to domain
func (c *DNS1Cloud) AddRecord(
	ctx context.Context,
	domainID uint64,
	record Record,
) (Record, error) {
	var (
		res Record
		cmd command
		err error
	)

	switch record.TypeRecord {
	case RecordTypeA:
		cmd, err = makeAddRecordACommand(domainID, record)
	case RecordTypeAAAA:
		cmd, err = makeAddRecordAAAACommand(domainID, record)
	case RecordTypeCNAME:
		cmd, err = makeAddRecordCNAMECommand(domainID, record)
	case RecordTypeMX:
		cmd, err = makeAddRecordMXCommand(domainID, record)
	case RecordTypeNS:
		cmd, err = makeAddRecordNSCommand(domainID, record)
	case RecordTypeSRV:
		cmd, err = makeAddRecordSRVCommand(domainID, record)
	case RecordTypeTXT:
		cmd, err = makeAddRecordTXTCommand(domainID, record)
	default:
		err = errors.Errorf("unknown record type: %d", record.TypeRecord)
	}
	if err != nil {
		return res, err
	}

	err = c.send(ctx, cmd, &res)
	return res, err
}

// addRecordAParams parameters for request for creating A and AAAA records
type addRecordAParams struct {
	DomainID string `json:"DomainId"`
	IP       string `json:"IP"`
	Name     string `json:"Name"`
	TTL      string `json:"TTL,omitempty"`
}

// addRecordCNAMEParams parameters for request for creating CNAME records
type addRecordCNAMEParams struct {
	DomainID     string `json:"DomainId"`
	Name         string `json:"Name"`
	MnemonicName string `json:"MnemonicName"`
	TTL          string `json:"TTL,omitempty"`
}

// addRecordMXParams parameters for request for creating MX records
type addRecordMXParams struct {
	DomainID string `json:"DomainId"`
	HostName string `json:"HostName"`
	Priority string `json:"Priority"`
	TTL      string `json:"TTL,omitempty"`
}

// addRecordNSParams parameters for request for creating NS records
type addRecordNSParams struct {
	DomainID string `json:"DomainId"`
	HostName string `json:"HostName"`
	Name     string `json:"Name"`
	TTL      string `json:"TTL,omitempty"`
}

// addRecordSRVParams parameters for request for creating SRV records
type addRecordSRVParams struct {
	DomainID string `json:"DomainId"`
	Service  string `json:"Service"`
	Proto    string `json:"Proto"`
	Name     string `json:"Name"`
	Priority string `json:"Priority"`
	Weight   string `json:"Weight"`
	Port     string `json:"Port"`
	Target   string `json:"Target"`
	TTL      string `json:"TTL,omitempty"`
}

// addRecordTXTParams parameters for request for creating TXT records
type addRecordTXTParams struct {
	DomainID string `json:"DomainId"`
	HostName string `json:"HostName"`
	Text     string `json:"Text"`
	TTL      string `json:"TTL,omitempty"`
}

func makeAddRecordACommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	ip := net.ParseIP(record.IP).To4()
	if ip == nil {
		return command{}, errors.Errorf("IP %q is incorrect", record.IP)
	}

	params := addRecordAParams{
		DomainID: strconv.FormatUint(domainID, 10),
		IP:       ip.String(),
		Name:     record.HostName,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recorda",
		params:   &params,
	}, nil
}

func makeAddRecordAAAACommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	ip := net.ParseIP(record.IP).To16()
	if ip == nil {
		return command{}, errors.Errorf("IP %q is incorrect", record.IP)
	}

	params := addRecordAParams{
		DomainID: strconv.FormatUint(domainID, 10),
		IP:       ip.String(),
		Name:     record.HostName,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recordaaaa",
		params:   &params,
	}, nil
}

func makeAddRecordCNAMECommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := addRecordCNAMEParams{
		DomainID:     strconv.FormatUint(domainID, 10),
		Name:         record.HostName,
		MnemonicName: record.MnemonicName,
		TTL:          ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recordcname",
		params:   &params,
	}, nil
}

func makeAddRecordMXCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := addRecordMXParams{
		DomainID: strconv.FormatUint(domainID, 10),
		HostName: record.HostName,
		Priority: record.Priority,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recordmx",
		params:   &params,
	}, nil
}

func makeAddRecordNSCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := addRecordNSParams{
		DomainID: strconv.FormatUint(domainID, 10),
		HostName: record.HostName,
		Name:     record.ExtHostName,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recordns",
		params:   &params,
	}, nil
}

func makeAddRecordSRVCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := addRecordSRVParams{
		DomainID: strconv.FormatUint(domainID, 10),
		Service:  record.Service,
		Proto:    record.Proto,
		Name:     record.HostName,
		Priority: record.Priority,
		Weight:   record.Weight,
		Port:     record.Port,
		Target:   record.Target,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recordsrv",
		params:   &params,
	}, nil
}

func makeAddRecordTXTCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := addRecordTXTParams{
		DomainID: strconv.FormatUint(domainID, 10),
		HostName: record.HostName,
		Text:     record.Text,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPost,
		endpoint: "dns/recordtxt",
		params:   &params,
	}, nil
}

func getTTL(ttl uint32) (string, error) {
	if ttl != 0 {
		str := strconv.FormatUint(uint64(ttl), 10)
		for _, t := range validTTLs {
			if t == str {
				return str, nil
			}
		}
		return "", errors.Errorf("TTL %q is not valid", str)
	}
	return "", nil
}

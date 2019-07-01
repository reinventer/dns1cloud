package dns1cloud

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// UpdateRecord updates record
func (c *DNS1Cloud) UpdateRecord(
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
		cmd, err = makeUpdateRecordACommand(domainID, record)
	case RecordTypeAAAA:
		cmd, err = makeUpdateRecordAAAACommand(domainID, record)
	case RecordTypeCNAME:
		cmd, err = makeUpdateRecordCNAMECommand(domainID, record)
	case RecordTypeMX:
		cmd, err = makeUpdateRecordMXCommand(domainID, record)
	case RecordTypeNS:
		cmd, err = makeUpdateRecordNSCommand(domainID, record)
	case RecordTypeSRV:
		cmd, err = makeUpdateRecordSRVCommand(domainID, record)
	case RecordTypeTXT:
		cmd, err = makeUpdateRecordTXTCommand(domainID, record)
	default:
		err = errors.Errorf("unknown record type: %d", record.TypeRecord)
	}
	if err != nil {
		return res, err
	}

	err = c.send(ctx, cmd, &res)
	return res, err
}

// updateRecordAParams parameters for request for updating A and AAAA records
type updateRecordAParams struct {
	DomainID string `json:"DomainId"`
	IP       string `json:"IP"`
	Name     string `json:"Name"`
	TTL      string `json:"TTL,omitempty"`
}

// updateRecordCNAMEParams parameters for request for updating CNAME records
type updateRecordCNAMEParams struct {
	DomainID     string `json:"DomainId"`
	Name         string `json:"Name"`
	MnemonicName string `json:"MnemonicName"`
	TTL          string `json:"TTL,omitempty"`
}

// updateRecordMXParams parameters for request for updating MX records
type updateRecordMXParams struct {
	DomainID string `json:"DomainId"`
	HostName string `json:"HostName"`
	Priority string `json:"Priority"`
	TTL      string `json:"TTL,omitempty"`
}

// updateRecordNSParams parameters for request for updating NS records
type updateRecordNSParams struct {
	DomainID string `json:"DomainId"`
	HostName string `json:"HostName"`
	Name     string `json:"Name"`
	TTL      string `json:"TTL,omitempty"`
}

// updateRecordSRVParams parameters for request for updating SRV records
type updateRecordSRVParams struct {
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

// updateRecordTXTParams parameters for request for updating TXT records
type updateRecordTXTParams struct {
	DomainID string `json:"DomainId"`
	Name     string `json:"Name"`
	Text     string `json:"Text"`
	TTL      string `json:"TTL,omitempty"`
}

func makeUpdateRecordACommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	ip := net.ParseIP(record.IP).To4()
	if ip == nil {
		return command{}, errors.Errorf("IP %q is incorrect", record.IP)
	}

	params := updateRecordAParams{
		DomainID: strconv.FormatUint(domainID, 10),
		IP:       ip.String(),
		Name:     record.HostName,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recorda/%d", record.ID),
		params:   &params,
	}, nil
}

func makeUpdateRecordAAAACommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	ip := net.ParseIP(record.IP).To16()
	if ip == nil {
		return command{}, errors.Errorf("IP %q is incorrect", record.IP)
	}

	params := updateRecordAParams{
		DomainID: strconv.FormatUint(domainID, 10),
		IP:       ip.String(),
		Name:     record.HostName,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recordaaaa/%d", record.ID),
		params:   &params,
	}, nil
}

func makeUpdateRecordCNAMECommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := updateRecordCNAMEParams{
		DomainID:     strconv.FormatUint(domainID, 10),
		Name:         record.HostName,
		MnemonicName: record.MnemonicName,
		TTL:          ttl,
	}

	return command{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recordcname/%d", record.ID),
		params:   &params,
	}, nil
}

func makeUpdateRecordMXCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := updateRecordMXParams{
		DomainID: strconv.FormatUint(domainID, 10),
		HostName: record.HostName,
		Priority: record.Priority,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recordmx/%d", record.ID),
		params:   &params,
	}, nil
}

func makeUpdateRecordNSCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := updateRecordNSParams{
		DomainID: strconv.FormatUint(domainID, 10),
		HostName: record.HostName,
		Name:     record.ExtHostName,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recordns/%d", record.ID),
		params:   &params,
	}, nil
}

func makeUpdateRecordSRVCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := updateRecordSRVParams{
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
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recordsrv/%d", record.ID),
		params:   &params,
	}, nil
}

func makeUpdateRecordTXTCommand(domainID uint64, record Record) (command, error) {
	ttl, err := getTTL(record.TTL)
	if err != nil {
		return command{}, err
	}

	params := updateRecordTXTParams{
		DomainID: strconv.FormatUint(domainID, 10),
		Name:     record.HostName,
		Text:     record.Text,
		TTL:      ttl,
	}

	return command{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("dns/recordtxt/%d", record.ID),
		params:   &params,
	}, nil
}

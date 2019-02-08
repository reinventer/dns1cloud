package dns1cloud

import (
	"fmt"
	"strings"
	"time"
)

// State is a state of the domain or record
type State uint8

const (
	StateNew State = iota
	StateActive
)

// UnmarshalJSON sets state of domain or record from JSON bytes
func (s *State) UnmarshalJSON(b []byte) error {
	state := string(b)
	switch state {
	case "null", `"New"`:
		*s = StateNew
	case `"Active"`:
		*s = StateActive
	default:
		return fmt.Errorf("unknown state %s", state)
	}
	return nil
}

// DateTime is a wrapper for unmarshaling time from JSON
type DateTime struct {
	time.Time
}

// UnmarshalJSON sets time from JSON bytes
func (d *DateTime) UnmarshalJSON(b []byte) error {
	var (
		str = string(b)
		t   time.Time
		err error
	)
	if str != "null" {
		t, err = time.Parse("2006-01-02T15:04:05.999999999", strings.Trim(string(b), `"`))
		if err != nil {
			return err
		}
	}

	d.Time = t
	return nil
}

// RecordType is a type of record
type RecordType uint8

const (
	// RecordTypeA is a type "A"
	RecordTypeA RecordType = iota
	// RecordTypeAAAA is a type "AAAA"
	RecordTypeAAAA
	// RecordTypeMX is a type "MX"
	RecordTypeMX
	// RecordTypeCNAME is a type "CNAME"
	RecordTypeCNAME
	// RecordTypeTXT is a type "TXT"
	RecordTypeTXT
	// RecordTypeNS is a type "NS"
	RecordTypeNS
	// RecordTypeSRV is a type "SRV"
	RecordTypeSRV
)

// UnmarshalJSON sets record type from JSON bytes
func (r *RecordType) UnmarshalJSON(b []byte) error {
	t := string(b)
	switch t {
	case `"A"`:
		*r = RecordTypeA
	case `"AAAA"`:
		*r = RecordTypeAAAA
	case `"MX"`:
		*r = RecordTypeMX
	case `"CNAME"`:
		*r = RecordTypeCNAME
	case `"TXT"`:
		*r = RecordTypeTXT
	case `"NS"`:
		*r = RecordTypeNS
	case `"SRV"`:
		*r = RecordTypeSRV
	default:
		return fmt.Errorf("unknown record type %s", t)
	}
	return nil
}

// Domain represents domain
type Domain struct {
	ID            uint64   `json:"ID"`
	Name          string   `json:"Name"`
	TechName      string   `json:"TechName"`
	State         State    `json:"State"`
	DateCreate    DateTime `json:"DateCreate"`
	IsDelegate    bool     `json:"IsDelegate"`
	LinkedRecords []Record `json:"LinkedRecords"`
}

// Domain represents record of domain
type Record struct {
	ID                   uint64     `json:"ID"`
	TypeRecord           RecordType `json:"TypeRecord"`
	IP                   string     `json:"IP"`
	HostName             string     `json:"HostName"`
	Priority             string     `json:"Priority"`
	Text                 string     `json:"Text"`
	MnemonicName         string     `json:"MnemonicName"`
	ExtHostName          string     `json:"ExtHostName"`
	Weight               string     `json:"Weight"`
	Port                 string     `json:"Port"`
	Target               string     `json:"Target"`
	Proto                string     `json:"Proto"`
	Service              string     `json:"Service"`
	TTL                  int        `json:"TTL"`
	State                State      `json:"State"`
	DateCreate           DateTime   `json:"DateCreate"`
	CanonicalDescription string     `json:"CanonicalDescription"`
}

package parser

import (
	"github.com/fossoreslp/go-dns/dns/record-names"
	"github.com/fossoreslp/go-dns/dns/record-types"
)

// Set is a set of DNS zones
type Set map[string]Zone

// Zone is a DNS zone
type Zone struct {
	Exclusive bool
	Entries   map[string]Entry
}

// Entry is one particular location in a DNS zone
type Entry map[names.TYPE][]record.Record

// GetRecordsOfType returns all records in an entry matching the specified type
func (e *Entry) GetRecordsOfType(t names.QTYPE) []record.Record {
	switch t {
	case names.AXFR, names.QTYPE_ANY:
		return nil
	case names.MAILB:
		return append((*e)[names.MD], (*e)[names.MF]...)
	case names.MAILA:
		out := append((*e)[names.MB], (*e)[names.MG]...)
		out = append(out, (*e)[names.MR]...)
		return append(out, (*e)[names.MINFO]...)
	default:
		return (*e)[names.TYPE(t)]
	}
}

// UnmarshalTOML is a function called by the TOML parser to properly decode the zones file entries.
func (e *Entry) UnmarshalTOML(decode func(interface{}) error) error {
	a := new(record.Records)
	if err := decode(a); err != nil {
		return err
	}
	o, err := a.Decode()
	if err != nil {
		return err
	}
	*e = Entry(o)
	return nil
}

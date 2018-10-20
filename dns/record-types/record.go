package record

import (
	"github.com/fossoreslp/go-dns/dns/record-names"
)

// Record an interface for all supported DNS record types
type Record interface {
	String() string
	Parse(string) error
	Encode() []byte
	Decode(message []byte, recordStart, recordLength int) error
	Type() names.TYPE
}

// Records is a struct used to decode the zones toml file.
type Records struct {
	A     []string
	AAAA  []string
	CAA   []string
	CNAME []string
	MX    []string
	NS    []string
	PTR   []string
	SOA   []string
	SRV   []string
	TXT   []string
}

// Decode returns the records in their proper individual formats
func (r Records) Decode() (map[names.TYPE][]Record, error) {
	rs := make(map[names.TYPE][]Record)

	rs[names.A] = make([]Record, 0)
	for _, s := range r.A {
		n := new(A)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.A] = append(rs[names.A], n)
	}

	rs[names.AAAA] = make([]Record, 0)
	for _, s := range r.AAAA {
		n := new(AAAA)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.AAAA] = append(rs[names.AAAA], n)
	}

	rs[names.CAA] = make([]Record, 0)
	for _, s := range r.CAA {
		n := new(CAA)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.CAA] = append(rs[names.CAA], n)
	}

	rs[names.CNAME] = make([]Record, 0)
	for _, s := range r.CNAME {
		n := new(CNAME)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.CNAME] = append(rs[names.CNAME], n)
	}

	rs[names.MX] = make([]Record, 0)
	for _, s := range r.MX {
		n := new(MX)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.MX] = append(rs[names.MX], n)
	}

	rs[names.NS] = make([]Record, 0)
	for _, s := range r.NS {
		n := new(NS)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.NS] = append(rs[names.NS], n)
	}

	rs[names.PTR] = make([]Record, 0)
	for _, s := range r.PTR {
		n := new(PTR)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.PTR] = append(rs[names.PTR], n)
	}

	rs[names.SOA] = make([]Record, 0)
	for _, s := range r.SOA {
		n := new(SOA)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.SOA] = append(rs[names.SOA], n)
	}

	rs[names.SRV] = make([]Record, 0)
	for _, s := range r.SRV {
		n := new(SRV)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.SRV] = append(rs[names.SRV], n)
	}

	rs[names.TXT] = make([]Record, 0)
	for _, s := range r.TXT {
		n := new(TXT)
		err := n.Parse(s)
		if err != nil {
			return nil, err
		}
		rs[names.TXT] = append(rs[names.TXT], n)
	}
	return rs, nil
}

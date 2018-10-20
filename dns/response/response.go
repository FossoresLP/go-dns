package response

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
	"github.com/fossoreslp/go-dns/dns/record-types"
)

// Response is used to store DNS records (the answers, ns and additional sections of all requests)
type Response struct {
	Name       label.Label
	Type       names.TYPE
	Class      names.CLASS
	TTL        uint32
	DataLength uint16
	Data       []byte
	Record     record.Record
}

func (r Response) String() string {
	return fmt.Sprintf("Record for %q, Type: %d, Class: %d, TTL: %d, Data length: %d, Data: %X", r.Name, r.Type, r.Class, r.TTL, r.DataLength, r.Data)
}

// Parse parses all records of a DNS message
func Parse(message []byte, start int, answers uint16, ns uint16, additional uint16) (Answers, NS, Additional []Response, err error) {
	Answers, start, err = parseSection(message, start, answers)
	if err != nil {
		return nil, nil, nil, err
	}
	NS, start, err = parseSection(message, start, ns)
	if err != nil {
		return nil, nil, nil, err
	}
	Additional, _, err = parseSection(message, start, additional)
	return
}

func parseSection(message []byte, start int, elements uint16) (out []Response, end int, err error) {
	end = start
	for ; elements > 0; elements-- {
		l, e, err := label.GetLabelsFromMessage(message, end)
		if err != nil {
			return nil, 0, err
		}
		if len(message) < e+10 {
			return nil, 0, errors.New("message to short for RR information")
		}
		t := binary.BigEndian.Uint16(message[e : e+2])
		c := binary.BigEndian.Uint16(message[e+2 : e+4])
		s := binary.BigEndian.Uint32(message[e+4 : e+8])
		r := binary.BigEndian.Uint16(message[e+8 : e+10])
		if len(message) < e+10+int(r) {
			return nil, 0, errors.New("data exceeds message lenght")
		}
		rt := getRecordType(names.TYPE(t))
		if rt != nil {
			err = rt.Decode(message, e+10, int(r))
			if err != nil {
				return nil, 0, err
			}
		}
		out = append(out, Response{l, names.TYPE(t), names.CLASS(c), s, r, message[e+10 : e+10+int(r)], rt})
		end = e + 10 + int(r)
	}
	return
}

func getRecordType(t names.TYPE) record.Record {
	switch t {
	case names.A:
		return new(record.A)
	case names.AAAA:
		return new(record.AAAA)
	case names.CAA:
		return new(record.CAA)
	case names.CNAME:
		return new(record.CNAME)
	case names.MX:
		return new(record.MX)
	case names.NS:
		return new(record.NS)
	case names.PTR:
		return new(record.PTR)
	case names.SOA:
		return new(record.SOA)
	case names.SRV:
		return new(record.SRV)
	case names.TXT:
		return new(record.TXT)
	default:
		return nil
	}
}

// New returns a new answer for the supplied contents
func New(name label.Label, t names.TYPE, ttl uint32, data []byte) Response {
	if len(data) > math.MaxUint16 {
		panic("data exceeds max record size")
	}
	return Response{name, t, names.IN, ttl, uint16(len(data)), data, nil}
}

// Encode returns the Response in DNS message format
func (r Response) Encode() []byte {
	r.DataLength = uint16(len(r.Data))
	b := r.Name.Encode()
	a := make([]byte, 10+int(r.DataLength))
	binary.BigEndian.PutUint16(a[:2], uint16(r.Type))
	binary.BigEndian.PutUint16(a[2:4], uint16(r.Class))
	binary.BigEndian.PutUint32(a[4:8], r.TTL)
	binary.BigEndian.PutUint16(a[8:10], r.DataLength)
	copy(a[10:], r.Data)
	return append(b, a...)
}

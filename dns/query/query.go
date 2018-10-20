package query

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
)

// Query represents a DNS query
type Query struct {
	Name  label.Label
	Type  names.QTYPE
	Class names.QCLASS
}

func (q Query) String() string {
	return fmt.Sprintf("Query for %q, Type: %d, Class: %d", q.Name, q.Type, q.Class)
}

// Parse returns an array containing all questions included in a message
func Parse(message []byte, number uint16) ([]Query, int, error) {
	queries := make([]Query, 0)
	if len(message) < 12 {
		return nil, 0, errors.New("message too short")
	}
	pos := 12
	for ; number > 0; number-- {
		labels, moved, err := label.GetLabelsFromMessage(message, pos)
		if err != nil {
			return nil, 0, err
		}
		pos += moved
		if len(message) < pos+4 {
			return nil, 0, errors.New("message too short for QTYPE and QCLASS")
		}
		t := binary.BigEndian.Uint16(message[pos : pos+2])
		class := binary.BigEndian.Uint16(message[pos+2 : pos+4])
		pos += 4
		queries = append(queries, Query{labels, names.QTYPE(t), names.QCLASS(class)})
	}
	return queries, pos, nil
}

// New returns a new Query for a record of type for name
func New(name label.Label, t names.QTYPE) Query {
	return Query{name, t, names.QCLASS(names.IN)}
}

// Encode returns the query in DNS message format
func (q Query) Encode() []byte {
	b := q.Name.Encode()
	a := make([]byte, 4)
	binary.BigEndian.PutUint16(a[:2], uint16(q.Type))
	binary.BigEndian.PutUint16(a[2:], uint16(q.Class))
	return append(b, a...)
}

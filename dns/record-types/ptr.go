package record

import (
	"errors"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
)

// PTR is the type for a PTR DNS record
type PTR struct {
	Label label.Label
}

// Type returns the record type
func (r PTR) Type() names.TYPE {
	return names.PTR
}

func (r PTR) String() string {
	return r.Label.String()
}

// Parse stores the input in the record
func (r *PTR) Parse(i string) error {
	l, err := label.Parse(i)
	if err != nil {
		return err
	}
	r.Label = l
	return nil
}

// Encode encodes the record to DNS message format
func (r PTR) Encode() []byte {
	return r.Label.Encode()
}

// Decode parses the input from DNS message format
func (r *PTR) Decode(i []byte, start, length int) error {
	l, end, err := label.GetLabelsFromMessage(i, start)
	if err != nil {
		return err
	}
	if end-start > length {
		return errors.New("Label exceeds data length")
	}
	r.Label = l
	return nil
}

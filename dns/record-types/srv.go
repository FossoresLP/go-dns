package record

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
)

// SRV represents a SRV DNS record
type SRV struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Host     label.Label
}

// Type returns the record type
func (r SRV) Type() names.TYPE {
	return names.SRV
}

func (r SRV) String() string {
	return fmt.Sprintf("%d %d %d %s", r.Priority, r.Weight, r.Port, r.Host.String())
}

// Parse stores the input in SRV
func (r *SRV) Parse(i string) error {
	s := strings.Split(i, " ")
	if len(s) != 4 {
		return errors.New("SRV record has to be in format \"Priority Weight Port Host\"")
	}
	priority, err := strconv.ParseUint(s[0], 10, 16)
	if err != nil {
		return err
	}
	weight, err := strconv.ParseUint(s[1], 10, 16)
	if err != nil {
		return err
	}
	port, err := strconv.ParseUint(s[2], 10, 16)
	if err != nil {
		return err
	}
	l, err := label.Parse(s[3])
	if err != nil {
		return err
	}
	r.Priority = uint16(priority)
	r.Weight = uint16(weight)
	r.Port = uint16(port)
	r.Host = l
	return nil
}

// Encode returns SRV in DNS message format
func (r SRV) Encode() []byte {
	out := make([]byte, 6)
	binary.BigEndian.PutUint16(out[:2], r.Priority)
	binary.BigEndian.PutUint16(out[2:4], r.Weight)
	binary.BigEndian.PutUint16(out[4:], r.Port)
	return append(out, r.Host.Encode()...)
}

// Decode extracts the SRV DNS record from a DNS message
func (r *SRV) Decode(i []byte, start, length int) error {
	if length < 6 {
		return errors.New("data too short for SRV record")
	}
	r.Priority = binary.BigEndian.Uint16(i[start : start+2])
	r.Weight = binary.BigEndian.Uint16(i[start+2 : start+4])
	r.Port = binary.BigEndian.Uint16(i[start+4 : start+6])
	l, end, err := label.GetLabelsFromMessage(i, start+6)
	if err != nil {
		return err
	}
	if end-start > length {
		return errors.New("label exceeds data length")
	}
	r.Host = l
	return nil
}

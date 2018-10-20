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

// MX is used to store MX DNS records
type MX struct {
	Priority uint16
	Name     label.Label
}

// Type returns the record type
func (r MX) Type() names.TYPE {
	return names.MX
}

func (r MX) String() string {
	return fmt.Sprintf("%s %d", r.Name.String(), r.Priority)
}

// Parse stores the input in the MX
func (r *MX) Parse(i string) error {
	e := strings.Split(i, " ")
	if len(e) != 2 {
		return errors.New("MX record needs to be in format mail.example.com 10")
	}
	l, err := label.Parse(e[0])
	if err != nil {
		return err
	}
	r.Name = l
	p, err := strconv.ParseUint(e[1], 10, 16)
	if err != nil {
		return errors.New("second section of MX record has to be uint16")
	}
	r.Priority = uint16(p)
	return nil
}

// Encode returns the record in DNS message format
func (r MX) Encode() []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, r.Priority)
	return append(out, r.Name.Encode()...)
}

// Decode stores the message data in the MX
func (r *MX) Decode(i []byte, start, length int) error {
	if length < 2 {
		return errors.New("data too short for MX record")
	}
	r.Priority = binary.BigEndian.Uint16(i[start : start+2])
	l, end, err := label.GetLabelsFromMessage(i, start+2)
	if err != nil {
		return err
	}
	if end-start > length {
		return errors.New("label exceeds data length")
	}
	r.Name = l
	return nil
}

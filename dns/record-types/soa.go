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

// SOA represents a SOA DNS record
type SOA struct {
	MName   label.Label
	RName   label.Label
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minimum uint32
}

// Type returns the record type
func (r SOA) Type() names.TYPE {
	return names.SOA
}

func (r SOA) String() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", r.MName.String(), r.RName.String(), r.Serial, r.Refresh, r.Retry, r.Expire, r.Minimum)
}

// Parse stores the input in SOA
func (r *SOA) Parse(i string) error {
	s := strings.Split(i, " ")
	if len(s) != 7 {
		return errors.New("SOA record has to be in format \"Primary Hostmaster Serial Refresh Retry Expire Minimum\"")
	}

	mname, err := label.Parse(s[0])
	if err != nil {
		return err
	}
	r.MName = mname

	rname, err := label.Parse(s[1])
	if err != nil {
		return err
	}
	r.RName = rname

	serial, err := strconv.ParseUint(s[2], 10, 32)
	if err != nil {
		return err
	}
	r.Serial = uint32(serial)

	refresh, err := strconv.ParseUint(s[3], 10, 32)
	if err != nil {
		return err
	}
	r.Refresh = uint32(refresh)

	retry, err := strconv.ParseUint(s[4], 10, 32)
	if err != nil {
		return err
	}
	r.Retry = uint32(retry)

	expire, err := strconv.ParseUint(s[5], 10, 32)
	if err != nil {
		return err
	}
	r.Expire = uint32(expire)

	minimum, err := strconv.ParseUint(s[6], 10, 32)
	if err != nil {
		return err
	}
	r.Minimum = uint32(minimum)

	return nil
}

// Encode returns the record in DNS message format
func (r SOA) Encode() []byte {
	s := append(r.MName.Encode(), r.RName.Encode()...)
	b := make([]byte, 20)
	binary.BigEndian.PutUint32(b[:4], r.Serial)
	binary.BigEndian.PutUint32(b[4:8], r.Refresh)
	binary.BigEndian.PutUint32(b[8:12], r.Retry)
	binary.BigEndian.PutUint32(b[12:16], r.Expire)
	binary.BigEndian.PutUint32(b[16:], r.Minimum)
	return append(s, b...)
}

// Decode parses the input from DNS message format
func (r *SOA) Decode(i []byte, start, lenght int) error {
	mname, m, err := label.GetLabelsFromMessage(i, start)
	if err != nil {
		return err
	}
	if m > lenght {
		return errors.New("label exceeds data lenght")
	}
	r.MName = mname

	rname, n, err := label.GetLabelsFromMessage(i, start+m)
	if err != nil {
		return err
	}
	if n > lenght {
		return errors.New("label exceeds data lenght")
	}
	r.RName = rname

	nums := i[start+m+n : start+lenght]
	if len(nums) != 20 {
		return errors.New("data lenght does not fit record")
	}

	r.Serial = binary.BigEndian.Uint32(nums[:4])
	r.Refresh = binary.BigEndian.Uint32(nums[4:8])
	r.Retry = binary.BigEndian.Uint32(nums[8:12])
	r.Expire = binary.BigEndian.Uint32(nums[12:16])
	r.Minimum = binary.BigEndian.Uint32(nums[16:])

	return nil
}

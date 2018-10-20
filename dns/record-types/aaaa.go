package record

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/fossoreslp/go-dns/dns/record-names"
)

// AAAA is used to store AAAA DNS records
type AAAA struct {
	IPv6 [16]byte
}

// Type returns the record type
func (r AAAA) Type() names.TYPE {
	return names.AAAA
}

func (r AAAA) String() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", r.IPv6[:2], r.IPv6[2:4], r.IPv6[4:6], r.IPv6[6:8], r.IPv6[8:10], r.IPv6[10:12], r.IPv6[12:14], r.IPv6[14:16])
}

// Parse stores an IPv6 string in the AAAA
func (r *AAAA) Parse(i string) error {
	s := strings.Split(i, ":")
	if len(s) != 8 {
		return errors.New("IPv6 address has to consist of 8 colon seperated blocks")
	}
	for i, v := range s {
		if len(v) != 4 {
			return errors.New("IPv6 block has to consist of 4 hex characters")
		}
		o1, err := hex.DecodeString(v[:2])
		if err != nil {
			return err
		}
		if len(o1) != 1 {
			return errors.New("HEX decode of 2 letters did not produce one byte")
		}
		o2, err := hex.DecodeString(v[2:4])
		if err != nil {
			return err
		}
		if len(o2) != 1 {
			return errors.New("HEX decode of 2 letters did not produce one byte")
		}
		r.IPv6[i*2] = o1[0]
		r.IPv6[i*2+1] = o2[0]
	}
	return nil
}

// Encode encodes the AAAA into DNS message format
func (r AAAA) Encode() []byte {
	return r.IPv6[:]
}

// Decode decodes an AAAA record from DNS message format into AAAA
func (r *AAAA) Decode(i []byte, start, length int) error {
	if length != 16 {
		return errors.New("invalid data length for AAAA record")
	}
	copy(r.IPv6[:], i[start:start+16])
	return nil
}

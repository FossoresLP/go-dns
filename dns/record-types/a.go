package record

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fossoreslp/go-dns/dns/record-names"
)

// The A type is used to store an A DNS record
type A struct {
	IPv4 [4]byte
}

// Type returns the record type
func (r A) Type() names.TYPE {
	return names.A
}

func (r A) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", r.IPv4[0], r.IPv4[1], r.IPv4[2], r.IPv4[3])
}

// Parse reads the string representation of an IPv4 into the A record
func (r *A) Parse(i string) error {
	s := strings.Split(i, ".")
	if len(s) != 4 {
		return errors.New("IPv4 string has to have 4 dot seperated sections")
	}
	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(s[i])
		if err != nil {
			return err
		}
		r.IPv4[i] = uint8(n)
	}
	return nil
}

// Encode returns the A in DNS message format
func (r A) Encode() []byte {
	return r.IPv4[:]
}

// Decode stores an A from DNS message format
func (r *A) Decode(i []byte, start, length int) error {
	if length != 4 {
		return errors.New("invalid data lenght for A record")
	}
	copy(r.IPv4[:], i[start:start+4])
	return nil
}

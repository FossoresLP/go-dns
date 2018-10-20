package record

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fossoreslp/go-dns/dns/record-names"
)

// CAA is used to store CAA DNS records
type CAA struct {
	Flags byte
	Tag   string
	Value string
}

// Type returns the record type
func (r CAA) Type() names.TYPE {
	return names.CAA
}

func (r CAA) String() string {
	return fmt.Sprintf("%d %s %s", r.Flags, r.Tag, r.Value)
}

// Parse stores input in CAA
func (r *CAA) Parse(i string) error {
	s := strings.SplitN(i, " ", 3)
	if len(s) != 3 {
		return errors.New("CAA records need to have the format \"Flags Tag Value\"")
	}
	f, err := strconv.ParseUint(s[0], 10, 8)
	if err != nil {
		return err
	}
	if f != 128 {
		return errors.New("only first bit of flag may be set")
	}
	r.Flags = uint8(f)
	if len(s[1]) > 255 {
		return errors.New("tag length may not exceed 255 characters")
	}
	r.Tag = s[1]
	r.Value = s[2]
	return nil
}

// Encode returns the CAA in DNS message format
func (r CAA) Encode() []byte {
	var out []byte
	out = append(out, r.Flags)
	out = append(out, uint8(len(r.Tag)))
	out = append(out, []byte(r.Tag)...)
	out = append(out, []byte(r.Value)...)
	return out
}

// Decode extracts the CAA from DNS message format
func (r *CAA) Decode(i []byte, start, length int) error {
	if length < 2 {
		return errors.New("data is not long enough to contain CAA record")
	}
	r.Flags = i[start]
	tlen := i[start+1]
	if length < int(tlen)+2 {
		return errors.New("tag lenght exceeds data lenght")
	}
	r.Tag = string(i[start+2 : start+2+int(tlen)])
	r.Value = string(i[start+2+int(tlen) : start+length])
	return nil
}

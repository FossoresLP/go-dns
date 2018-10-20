package record

import "github.com/fossoreslp/go-dns/dns/record-names"

// TXT is used to store TXT DNS records
type TXT struct {
	Content string
}

// Type returns the record type
func (r TXT) Type() names.TYPE {
	return names.TXT
}

func (r TXT) String() string {
	return r.Content
}

// Parse stores the input in TXT
func (r *TXT) Parse(i string) error {
	r.Content = i
	return nil
}

// Encode returns the TXT in DNS message format
func (r TXT) Encode() []byte {
	return []byte(r.Content)
}

// Decode extracts the TXT DNS record from DNS message format
func (r *TXT) Decode(i []byte, start, length int) error {
	if length > 0 {
		r.Content = string(i[start : start+length])
	}
	return nil
}

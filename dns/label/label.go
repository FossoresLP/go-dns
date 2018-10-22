package label

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var re *regexp.Regexp = regexp.MustCompile(`^[a-z][a-z0-9]+$`)

// Label is a type used to store a DNS label
type Label []string

func (l Label) String() string {
	s := ""
	for _, v := range l {
		s += v + "."
	}
	return s
}

// GetLabelsFromMessage extracts a label from a message starting at start and returns that label as well as it's end position.
func GetLabelsFromMessage(data []byte, start int) (Label, int, error) {
	if start >= len(data) {
		return nil, 0, errors.New("label outside data range")
	}
	w := data[start:]
	labels := make([]string, 0)
	m := 0
	for {
		/*if len(w) < 1 { // Check not necessary as a slice defined using s[0:] can never be shorter than 1
			return nil, 0, errors.New("label section too short")
		}*/
		l := w[0]
		if l == 0 {
			m++
			break
		} else if l == 192 {
			if len(w) < 2 {
				return nil, 0, errors.New("label section to short for redirect")
			}
			s, _, err := GetLabelsFromMessage(data, int(w[1]))
			if err != nil {
				return nil, 0, err
			}
			return s, start + 2, nil
		} else if l >= 64 {
			return nil, 0, errors.New("labels cannot be larger than 63 bytes")
		} else if int(l) >= len(w) {
			return nil, 0, errors.New("label section lenght exceeds data lenght")
		}
		labels = append(labels, string(w[1:1+l]))
		w = w[1+l:]
		m += int(1 + l)
	}
	return labels, m, nil
}

// Encode encodes a label to the DNS message format
func (l Label) Encode() []byte {
	var out []byte
	for _, s := range l {
		out = append(out, uint8(len(s))) // TODO: make sure len cannot exceed uint8 and section lenght cannot exceed 63 bytes
		out = append(out, s...)
	}
	out = append(out, 0)
	return out
}

// Parse turns a string into a label or returns an error if that's not possible
func Parse(i string) (Label, error) {
	s := strings.Split(i, ".")
	if s[len(s)-1] == "" {
		s = s[:len(s)-1]
	}
	if len(s) < 1 {
		return nil, errors.New("label cannot be empty")
	}
	for n := range s {
		s[n] = strings.ToLower(s[n])
		fmt.Println(s[n])
		fmt.Println(re.FindString(s[n]))
		if !re.MatchString(s[n]) {
			return nil, errors.New("label has to start with a letter and may only contain letters and numbers")
		}
	}
	return s, nil
}

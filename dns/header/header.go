package header

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// Header is the header of a DNS message
type Header struct {
	ID              [2]byte
	Flags           [2]byte
	QuestionCount   uint16
	AnswerCount     uint16
	NSCount         uint16
	AdditionalCount uint16
}

func (h Header) String() string {
	return fmt.Sprintf("Message with ID %X.\nThe flags have the following values: QR: %t, OpCode: %d, AA: %t, TC: %t, RD: %t, RA: %t, Zero bits are actually zero: %t, RCode: %d.\nThe message contains %d questions, %d answers, %d authority pointers and %d additional records.", h.ID, h.IsResponse(), h.OpCode(), h.AuthoritativeAnswer(), h.Truncated(), h.RecursionDesired(), h.RecursionAvailable(), h.ZeroBits(), h.ResponseCode(), h.QuestionCount, h.AnswerCount, h.NSCount, h.AdditionalCount)
}

// Parse extracts the header from a DNS message
func Parse(in []byte) (*Header, error) {
	if len(in) < 12 {
		return nil, errors.New("message too short for header")
	}
	h := Header{}
	copy(h.ID[:], in[:2])
	copy(h.Flags[:], in[2:4])
	h.QuestionCount = binary.BigEndian.Uint16(in[4:6])
	h.AnswerCount = binary.BigEndian.Uint16(in[6:8])
	h.NSCount = binary.BigEndian.Uint16(in[8:10])
	h.AdditionalCount = binary.BigEndian.Uint16(in[10:12])
	return &h, nil
}

// IsResponse returns true if the QR bit is set
func (h Header) IsResponse() bool {
	return (h.Flags[0] >> 7) == 1
}

// OpCode returns the OpCode encoded in the flags
func (h Header) OpCode() uint8 {
	return (h.Flags[0] & 0x7F) >> 3
}

// AuthoritativeAnswer returns true if the AA bit is set
func (h Header) AuthoritativeAnswer() bool {
	return ((h.Flags[0] & 0x7) >> 2) == 1
}

// Truncated returns true if the TC bit is set
func (h Header) Truncated() bool {
	return ((h.Flags[0] & 0x3) >> 1) == 1
}

// RecursionDesired returns true if the RD bit is set
func (h Header) RecursionDesired() bool {
	return (h.Flags[0] & 0x1) == 1
}

// RecursionAvailable returns true if the RA bit is set
func (h Header) RecursionAvailable() bool {
	return (h.Flags[1] >> 7) == 1
}

// ZeroBits returns true if the Z bits are not set as required by the standard
func (h Header) ZeroBits() bool {
	return ((h.Flags[1] & 0x7F) >> 4) == 0
}

// ResponseCode returns the RCode encoded in the flags
func (h Header) ResponseCode() uint8 {
	return h.Flags[1] & 0xF
}

// NewQueryHeader returns a new Header with the settings for a standard request
func NewQueryHeader(rd bool) *Header {
	h := new(Header)
	_, err := rand.Read(h.ID[:])
	if err != nil {
		binary.BigEndian.PutUint16(h.ID[:], uint16(time.Now().UnixNano()))
	}
	if rd {
		binary.BigEndian.PutUint16(h.Flags[:], 0x0100) // QR: 0, OpCode: 0000, AA: 0, TC: 0, RD: 1, RA: 0, Z: 000, RCode: 0000
	} else {
		binary.BigEndian.PutUint16(h.Flags[:], 0x0000) // Everything is set to zero.
	}
	return h
}

// NewErrorHeader creates a new header used to convey an error
func NewErrorHeader(id [2]byte, aa bool, rcode uint8) *Header {
	if rcode > 10 { // 15 would be the theoretical limit for 4 bits but only the RCodes from 0 to 10 are defined
		return nil
	}
	f1 := byte(0x80)         // 1000 0000
	f2 := byte(0x80) | rcode // 1000 0000
	if aa {
		f1 = f1 | 0x04 // 0000 0100
	}

	h := new(Header)
	h.ID = id
	h.Flags[0] = f1
	h.Flags[1] = f2

	return h
}

// NewAnswerHeader returns a new header for a DNS response message
func NewAnswerHeader(id [2]byte, aa bool, rd bool) *Header {
	f1 := byte(0x80)
	f2 := byte(0x80)
	if aa {
		f1 = f1 | 0x04
	}
	if rd {
		f1 = f1 | 0x01
	}

	h := new(Header)
	h.ID = id
	h.Flags[0] = f1
	h.Flags[1] = f2

	return h
}

// Encode returns the header in DNS message format
func (h Header) Encode() []byte {
	b := make([]byte, 12)
	copy(b[:2], h.ID[:])
	copy(b[2:4], h.Flags[:])
	binary.BigEndian.PutUint16(b[4:6], h.QuestionCount)
	binary.BigEndian.PutUint16(b[6:8], h.AnswerCount)
	binary.BigEndian.PutUint16(b[8:10], h.NSCount)
	binary.BigEndian.PutUint16(b[10:], h.AdditionalCount)
	return b
}

package dnserror

import (
	"fmt"

	"github.com/fossoreslp/go-dns/dns/header"
	"github.com/fossoreslp/go-dns/dns/message"
	"github.com/fossoreslp/go-dns/dns/query"
)

// Error is the DNS error type provided by this package
type Error struct {
	RCode uint8
	AA    bool
}

func (e Error) Error() string {
	return fmt.Sprintf("DNS error with code %d. Authoritative response: %t", e.RCode, e.AA)
}

// Message returns the error as a DNS message
func (e Error) Message(id [2]byte, q ...query.Query) *message.Message {
	return message.New(header.NewErrorHeader(id, e.AA, e.RCode), q, nil, nil, nil)
}

// IsError returns true if the RCode is not 0 (NoError)
func (e Error) IsError() bool {
	return e.RCode != NoError
}

// New creates a new DNS error using the supplied data
func New(rcode uint8, aa bool) Error {
	return Error{rcode, aa}
}

// Success returns an Error with RCode 0 (NoError)
func Success() Error {
	return Error{NoError, false}
}

const (
	// No Error = 0
	NoError uint8 = iota
	// Format Error = 1
	FormatError
	// Server Failure = 2
	ServerFailure
	// Name Error = 3
	NameError
	// Not Implemented = 4
	NotImplemented
	// Refused = 5
	Refused
	// YX Domain = 6
	YXDomain
	// YX RR Set = 7
	YXRRSet
	// NX RR Set = 8
	NXRRSet
	// Not Auth = 9
	NotAuth
	// Not Zone = 10
	NotZone
)

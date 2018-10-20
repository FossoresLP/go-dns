package message

import (
	"errors"

	"github.com/fossoreslp/go-dns/dns/header"
	"github.com/fossoreslp/go-dns/dns/query"
	"github.com/fossoreslp/go-dns/dns/response"
)

// Message is used to store a DNS message
type Message struct {
	Header      *header.Header
	Questions   []query.Query
	Answers     []response.Response
	Authorities []response.Response
	Additional  []response.Response
}

// Parse parses a DNS message
func Parse(msg []byte) (*Message, error) {
	h, err := header.Parse(msg)
	if err != nil {
		return nil, err
	}
	if h.Truncated() {
		return nil, errors.New("truncated messages cannot be handled, yet")
	}
	q, qend, err := query.Parse(msg, h.QuestionCount)
	if err != nil {
		return nil, err
	}
	ans, auth, add, err := response.Parse(msg, qend, h.AnswerCount, h.NSCount, h.AdditionalCount)
	if err != nil {
		return nil, err
	}
	return &Message{h, q, ans, auth, add}, nil
}

func (msg Message) String() string {
	var out string
	out += msg.Header.String() + "\n"
	out += qArrToString(msg.Questions)
	out += rArrToString(msg.Answers)
	out += rArrToString(msg.Authorities)
	out += rArrToString(msg.Additional)
	return out
}

func qArrToString(qs []query.Query) string {
	var out string
	for _, q := range qs {
		out += q.String() + "\n"
	}
	return out
}

func rArrToString(rs []response.Response) string {
	var out string
	for _, r := range rs {
		out += r.String() + "\n"
	}
	return out
}

// New creates a new message with the supplied contents
func New(h *header.Header, qs []query.Query, as []response.Response, ns []response.Response, add []response.Response) *Message {
	if h == nil {
		return nil
	}
	h.QuestionCount = uint16(len(qs))
	h.AnswerCount = uint16(len(as))
	h.NSCount = uint16(len(ns))
	h.AdditionalCount = uint16(len(add))
	return &Message{h, qs, as, ns, add}
}

// Encode returns the message in encoded DNS message format
func (msg Message) Encode() []byte {
	msg.Header.QuestionCount = uint16(len(msg.Questions))

	msg.Header.AdditionalCount = uint16(len(msg.Answers))

	msg.Header.NSCount = uint16(len(msg.Authorities))

	msg.Header.AdditionalCount = uint16(len(msg.Additional))

	b := msg.Header.Encode()

	for _, v := range msg.Questions {
		b = append(b, v.Encode()...)
	}

	for _, v := range msg.Answers {
		b = append(b, v.Encode()...)
	}

	for _, v := range msg.Authorities {
		b = append(b, v.Encode()...)
	}

	for _, v := range msg.Additional {
		b = append(b, v.Encode()...)
	}

	return b
}

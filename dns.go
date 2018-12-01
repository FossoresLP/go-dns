package main

import (
	"fmt"
	"net"

	"github.com/fossoreslp/go-dns/dns/cache"
	"github.com/fossoreslp/go-dns/dns/error"
	"github.com/fossoreslp/go-dns/dns/header"
	"github.com/fossoreslp/go-dns/dns/message"
	"github.com/fossoreslp/go-dns/dns/passthrough"
	"github.com/fossoreslp/go-dns/dns/query"
	"github.com/fossoreslp/go-dns/dns/record-names"
	"github.com/fossoreslp/go-dns/dns/record-parser"
	"github.com/fossoreslp/go-dns/dns/response"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 53, IP: nil})
	if err != nil {
		panic(err)
	}
	defer listener.Close() //nolint: errcheck
	println("Listening...")

	set := parser.ParseZonesFile()
	println("Initialization finished")
	for {
		var buffer [512]byte
		rlen, remote, err := listener.ReadFromUDP(buffer[:])
		if err != nil {
			continue // TODO: Implement error logging
		}
		req, err := message.Parse(buffer[:rlen])
		if err != nil {
			continue // TODO: Implement error logging
		}
		var out []byte
		if !req.Header.IsResponse() && req.Header.QuestionCount > 0 {
			responses := make([]response.Response, 0)
			local := false
			for _, q := range req.Questions {
				if q.Class != names.QCLASS(names.IN) {
					out = dnserror.New(dnserror.NotImplemented, false).Message(req.Header.ID, q).Encode()
					break
				}
				resp, dnserr := findInLocalZones(q, set, req)
				if dnserr.IsError() {
					out = dnserr.Message(req.Header.ID, q).Encode()
					break
				}
				if resp != nil {
					responses = append(responses, resp...)
					local = true
					continue
				}
				resp = cache.GetRecords(q.Name, q.Type)
				if resp != nil {
					responses = append(responses, resp...)
					println("Cache hit")
					continue
				}
				println("Cache miss")
				resp, dnserr = passthrough.Resolve(q)
				if dnserr.IsError() {
					out = dnserr.Message(req.Header.ID, q).Encode()
					break
				}
				cache.Cache(resp)
				responses = append(responses, resp...)
			}
			if out == nil {
				out = message.New(header.NewAnswerHeader(req.Header.ID, local, req.Header.RecursionDesired()), req.Questions, responses, nil, nil).Encode()
			}
		} else {
			out = dnserror.New(dnserror.FormatError, false).Message(req.Header.ID).Encode()
		}

		_, err = listener.WriteToUDP(out, remote)
		if err != nil {
			fmt.Println("Failed to send response:", err.Error(), "- retrying")
			_, err = listener.WriteToUDP(out, remote)
			if err != nil {
				fmt.Println("Could not send response on second attempt:", err.Error(), "- giving up")
			} else {
				fmt.Println("Retry successful - please inspect first error message")
			}
		}
	}
}

func findInLocalZones(q query.Query, set *parser.Set, req *message.Message) ([]response.Response, dnserror.Error) {
	if set == nil {
		return nil, dnserror.Success()
	}
	e, excl := parser.Match(q.Name, set)
	if e == nil && excl {
		return nil, dnserror.New(dnserror.NameError, true)
	}
	if e == nil && !excl {
		return nil, dnserror.Success()
	}
	rs := e.GetRecordsOfType(q.Type)
	responses := make([]response.Response, 0)
	for _, r := range rs {
		responses = append(responses, response.New(q.Name, r.Type(), 60, r.Encode()))
	}
	return responses, dnserror.Success()
}

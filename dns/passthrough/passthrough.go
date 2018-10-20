package passthrough

import (
	"fmt"
	"net"
	"os"

	"github.com/fossoreslp/go-dns/dns/error"
	"github.com/fossoreslp/go-dns/dns/header"
	"github.com/fossoreslp/go-dns/dns/message"
	"github.com/fossoreslp/go-dns/dns/query"
	"github.com/fossoreslp/go-dns/dns/response"
)

var cResolveRequest chan query.Query
var cResolveResponse chan message.Message

func init() {
	cResolveRequest = make(chan query.Query)
	cResolveResponse = make(chan message.Message)
	go routine()
}

// Routine is the function used as a goroutine to resolve a DNS request with CloudFlare
func routine() {
	cfdns, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.ParseIP("1.1.1.1")})
	if err != nil {
		panic(err)
	}
	defer cfdns.Close() //nolint: errcheck
	var cfbuf [512]byte
	for {
		data := <-cResolveRequest
		msg := message.New(header.NewQueryHeader(true), []query.Query{data}, nil, nil, nil)
		_, err := cfdns.Write(msg.Encode())
		if err != nil {
			fmt.Println("Write failed with error:", err.Error(), "- trying to reconnect")
			c, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.ParseIP("1.1.1.1")})
			if err != nil {
				fmt.Println("Reconnect failed with error:", err.Error(), "- exiting")
				os.Exit(-1)
			}
			_, err = c.Write(msg.Encode())
			if err != nil {
				fmt.Println("Reconnect successful but write still failed with error:", err.Error(), "- exiting")
				os.Exit(-1)
			}
			fmt.Println("Reconnected successfully - replacing connection with newly established one")
			cfdns = c
		}

		n, err := cfdns.Read(cfbuf[:])
		if err != nil {
			fmt.Println("Read failed with error:", err.Error(), "- trying to reconnect")
			c, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 53, IP: net.ParseIP("1.1.1.1")})
			if err != nil {
				fmt.Println("Reconnect failed with error:", err.Error(), "- exiting")
				os.Exit(-1)
			}
			n, err = c.Read(cfbuf[:])
			if err != nil {
				fmt.Println("Reconnect successful but read still failed with error:", err.Error(), "- exiting")
				os.Exit(-1)
			}
			fmt.Println("Reconnected successfully - replacing connection with newly established one")
			cfdns = c
		}
		m, err := message.Parse(cfbuf[:n])
		if err != nil {
			fmt.Println("Failed to parse message:", err)
			cResolveResponse <- message.Message{}
			continue
		}
		cResolveResponse <- *m
	}
}

// Resolve resolves the query with the external DNS provider
func Resolve(q query.Query) ([]response.Response, dnserror.Error) {
	cResolveRequest <- q
	r := <-cResolveResponse
	if r.Header == nil {
		return nil, dnserror.New(dnserror.ServerFailure, false)
	}
	if r.Header.ResponseCode() != 0 {
		return nil, dnserror.New(r.Header.ResponseCode(), r.Header.AuthoritativeAnswer())
	}
	return r.Answers, dnserror.Success()
}

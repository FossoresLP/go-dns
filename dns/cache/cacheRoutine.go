package cache

import (
	"time"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
	"github.com/fossoreslp/go-dns/dns/record-types"
	"github.com/fossoreslp/go-dns/dns/response"
)

type recordTransfer struct {
	Record []recordStorage
	Type   names.TYPE
	Domain string
}

type recordStorage struct {
	Record    record.Record
	TTL       uint32
	StoreTime int64
}

var cRequest chan recordTransfer
var cResolved chan recordTransfer

func init() {
	cRequest = make(chan recordTransfer)
	cResolved = make(chan recordTransfer)
	go routine()
}

// Routine is the go routine used to cache records
func routine() {
	dataStore := make(map[names.TYPE]map[string][]recordStorage)
	for {
		req := <-cRequest
		if req.Record == nil {
			if t, ok := dataStore[req.Type]; ok {
				if r, ok := t[req.Domain]; ok {
					cResolved <- recordTransfer{r, req.Type, req.Domain}
					continue
				}
			}
			cResolved <- req
			continue
		}
		if _, ok := dataStore[req.Type]; ok {
			if _, ok := dataStore[req.Type][req.Domain]; ok {
				dataStore[req.Type][req.Domain] = append(dataStore[req.Type][req.Domain], req.Record...)
				continue
			}
			dataStore[req.Type][req.Domain] = req.Record
			continue
		}
		d := make(map[string][]recordStorage)
		d[req.Domain] = req.Record
		dataStore[req.Type] = d
	}
}

// GetRecords gets all records of a specific type and for a specific label
func GetRecords(lbl label.Label, t names.QTYPE) (out []response.Response) {
	switch t {
	case names.AXFR, names.QTYPE_ANY: // AXFR is only supported by authoritative nameservers and ANY will be deprecated soon
		return nil
	case names.MAILB: // Should return MD and MF
		return nil // These record types are not used and their implementation is therefore low priority
	case names.MAILA: // Should return MB, MG, MR and MINFO
		return nil // These record types are not used and their implementation is therefore low priority
	}
	cRequest <- recordTransfer{nil, names.TYPE(t), lbl.String()}
	res := <-cResolved
	for _, r := range res.Record {
		remaining := int64(r.TTL) - (time.Now().Unix() - r.StoreTime)
		if remaining < 1 {
			continue
		}
		d := r.Record.Encode()
		out = append(out, response.Response{Name: lbl, Type: res.Type, Class: names.IN, TTL: uint32(remaining), DataLength: uint16(len(d)), Data: d, Record: r.Record})
	}
	return
}

// Store takes a slice of DNS responses and adds them to the cache
func Store(res []response.Response) {
	for _, r := range res {
		cRequest <- recordTransfer{[]recordStorage{recordStorage{r.Record, r.TTL, time.Now().Unix()}}, r.Type, r.Name.String()}
	}
}

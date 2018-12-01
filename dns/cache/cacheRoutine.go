package cache

import (
	"time"

	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
	"github.com/fossoreslp/go-dns/dns/response"
)

type tCacheTransfer struct {
	lbl     label.Label
	t       names.TYPE
	records []RecordWrapper
	remove  bool
}

var cStore chan tCacheTransfer

var vStore *Node

func init() {
	cStore = make(chan tCacheTransfer)
	vStore = NewNode(NewStore())
	go routine()
}

// Routine is the go routine used to cache records
func routine() {
	for {
		t := <-cStore
		var node = vStore
		for i := len(t.lbl) - 1; i >= 0; i-- {
			if n := node.GetChild(t.lbl[i]); n != nil {
				node = n
				continue
			}
			n := NewNode(NewStore())
			err := node.AddChild(t.lbl[i], n) // No need to handle error as we just checked if the node already exists
			if err != nil {
				panic("Cache changed during synchronous operation.")
			}
			node = n
		}
		if t.remove {
			node.Content.(*Store).RemoveElement(t.t)
		} else {
			node.Content.(*Store).AddElement(t.t, t.records)
		}
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

	node := vStore
	for i := len(lbl) - 1; i >= 0; i-- {
		n := node.GetChild(lbl[i])
		if n == nil {
			return nil
		}
		node = n
	}

	records := node.Content.(*Store).GetElement(names.TYPE(t)) // No need to handle errors when checking for zero values

	if records == nil {
		return nil
	}

	for _, r := range records {
		remaining := int64(r.TTL) - (time.Now().Unix() - r.StoredAt)
		if remaining < 1 {
			cStore <- tCacheTransfer{r.Label, r.Record.Type(), nil, true}
			return nil
		}
		d := r.Record.Encode()
		out = append(out, response.Response{Name: lbl, Type: r.Record.Type(), Class: names.IN, TTL: uint32(remaining), DataLength: uint16(len(d)), Data: d, Record: r.Record})
	}
	return
}

// Cache takes a slice of DNS responses and adds them to the cache
func Cache(res []response.Response) {
	labels := make(map[string]map[names.TYPE][]RecordWrapper)
	for _, r := range res {
		lbl := r.Name.String()
		types, labelExists := labels[lbl]
		if !labelExists {
			labels[lbl] = make(map[names.TYPE][]RecordWrapper)
		}
		t, typeExists := types[r.Type]
		if !typeExists {
			labels[lbl][r.Type] = []RecordWrapper{RecordWrapper{r.Name, r.Record, r.TTL, time.Now().Unix()}}
			continue
		}
		labels[lbl][r.Type] = append(t, RecordWrapper{r.Name, r.Record, r.TTL, time.Now().Unix()})
	}
	for _, l := range labels {
		for typ, t := range l {
			cStore <- tCacheTransfer{t[0].Label, typ, t, false}
		}
	}
}

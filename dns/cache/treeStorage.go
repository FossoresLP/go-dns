package cache

import (
	"github.com/fossoreslp/go-dns/dns/label"
	"github.com/fossoreslp/go-dns/dns/record-names"
	"github.com/fossoreslp/go-dns/dns/record-types"
)

// NewNode returns a new node containing store
func NewNode(s *Store) *Node {
	return &Node{make(map[string]*Node), s}
}

// RecordWrapper is used to add information like TTL and cache time to a DNS record
type RecordWrapper struct {
	Label    label.Label
	Record   record.Record
	TTL      uint32
	StoredAt int64
}

// Store is a TreeStore used to store DNS records in a map by type
type Store struct {
	records map[names.TYPE][]RecordWrapper
}

// NewStore returns a new, initialized store
func NewStore() *Store {
	return &Store{make(map[names.TYPE][]RecordWrapper)}
}

// GetElement gets all records of a specific type stored by a node
func (s Store) GetElement(t names.TYPE) []RecordWrapper {
	if rs, ok := s.records[t]; ok {
		return rs
	}
	return nil
}

// AddElement adds a slice of records of a specific type to a node
func (s *Store) AddElement(t names.TYPE, rs []RecordWrapper) {
	if _, ok := s.records[t]; !ok {
		s.records[t] = rs
		return
	}
	s.records[t] = append(s.records[t], rs...)
}

// RemoveElement deletes all records of a specific type from a node
func (s *Store) RemoveElement(t names.TYPE) {
	if _, ok := s.records[t]; !ok {
		return //In case there are no records of the type, just ignore that
	}
	delete(s.records, t)
}

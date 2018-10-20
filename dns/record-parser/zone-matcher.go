package parser

import (
	"github.com/fossoreslp/go-dns/dns/label"
)

// FindMatchingZone tries to find the highest (if any) matching zone for label in set
func FindMatchingZone(l label.Label, set *Set) (*Zone, int) {
	var highestMatchingZone *Zone
	var sectionCountOfMatch int
	for s := 1; s < len(l); s++ {
		if val, ok := (*set)[concat(l[s:])]; ok {
			highestMatchingZone = &val
			sectionCountOfMatch = s
		}
	}
	return highestMatchingZone, sectionCountOfMatch
}

// FindMatchingEntry tries to find a matching entry for label in zone which occupies zoneSections sections of zone
func FindMatchingEntry(l label.Label, zone *Zone, zoneSections int) *Entry {
	if val, ok := (*zone).Entries[concat(l[:zoneSections])]; ok {
		return &val
	}
	return nil
}

// Match tries to find a matching entry for label in set and returns that entry and if the zone is exclusive
func Match(l label.Label, set *Set) (*Entry, bool) {
	z, c := FindMatchingZone(l, set)
	if z == nil {
		return nil, false
	}
	return FindMatchingEntry(l, z, c), z.Exclusive
}

func concat(s []string) (out string) {
	if len(s) == 1 {
		return s[0]
	}
	for _, p := range s {
		out += "." + p
	}
	out = out[1:]
	return
}

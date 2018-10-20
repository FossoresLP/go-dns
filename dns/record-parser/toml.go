package parser

import (
	"fmt"
	"os"

	"github.com/naoina/toml"
)

// ParseZonesFile parses the zone file
func ParseZonesFile() *Set {
	f, err := os.Open("zones.toml")
	if err != nil {
		return nil
	}
	defer f.Close() //nolint: errcheck
	set := new(Set)
	if err := toml.NewDecoder(f).Decode(set); err != nil {
		fmt.Printf("Failed to parse zones.toml: %s.\nContinuing without local zones.\n", err.Error())
		return nil
	}
	return set
}

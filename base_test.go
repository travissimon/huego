package huego

import (
	"fmt"
	"testing"
)

func TestDiscoverBase(t *testing.T) {
	bases, err := DiscoverBases()

	if err != nil {
		fmt.Errorf("Error discovering base: %v", err)
	}

	if err == nil {
		fmt.Errorf("Base is nil")
	}

	if len(bases) != 1 {
		fmt.Errorf("Expected 1 base, got: %v", len(bases))
	}

	fmt.Printf("Base ipaddress: %v\n", bases[0].InternalIp)
}

package huego

import (
	"testing"
)

func TestNothing(t *testing.T) {
	t.Pass()
}

/*
var testBase = getTestBase()

func getTestBase() *Base {
	bases, _ := DiscoverBases()

	testBase := &bases[0]
	return testBase
}

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
	fmt.Printf("Username: %v\n", bases[0].Username)
}

func TestLightNames(t *testing.T) {
	base := getTestBase()

	lightNames, _ := base.GetLightNames()
	for _, name := range lightNames {
		fmt.Printf("%v: %v\n", name.Id, name.Name)
	}
}

func TestGetLights(t *testing.T) {
	base := getTestBase()

	lights, _ := base.GetLights()
	fmt.Printf("len(lights): %v\n", len(lights))
	for _, light := range lights {
		fmt.Printf("%v: %v\n", light.Id, light)
	}
}

func TestChangeLight(t *testing.T) {
	base := getTestBase()

	lights, _ := base.GetLights()
	light := lights[0]

	light.State.Hue = 0
	light.SetStateWithTransition(10)
}

*/

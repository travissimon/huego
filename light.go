package Huego

type State struct {
	On        bool
	Hue       uint16
	Bri       uint8
	Sat       uint8
	X         float
	Y         float
	Ct        uint16
	Alert     string
	Effect    string
	Colormode string
	Reachable bool
}

type Light struct {
	prevState       *State
	State           *State
	Type            string
	Name            string
	Model           string
	SoftwareVersion string
}

func NewLight() *Light {
	light := new(Light)
	light.prevState = new(State)
	light.State = new(State)
	return light
}

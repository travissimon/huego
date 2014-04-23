package huego

import (
	"bytes"
	"fmt"
	"math"
)

type State struct {
	On        bool
	Hue       uint16
	Bri       uint8
	Sat       uint8
	X         float32
	Y         float32
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
	Name            *string
	Model           string
	SoftwareVersion string
}

func NewLight() Light {
	name := ""
	light := Light{
		new(State),
		new(State),
		"",
		&name,
		"",
		"",
	}
	return light
}

func (l *Light) ResetState() {
	l.prevState.On = l.State.On
	l.prevState.Hue = l.State.Hue
	l.prevState.Bri = l.State.Bri
	l.prevState.Sat = l.State.Sat
	l.prevState.X = l.State.X
	l.prevState.Y = l.State.Y
	l.prevState.Ct = l.State.Ct
	l.prevState.Alert = l.State.Alert
	l.prevState.Effect = l.State.Effect
	l.prevState.Colormode = l.State.Colormode
	l.prevState.Reachable = l.State.Reachable
}

// Returns the number of 'edit steps' between old and new hsv colour
// Used to prioritise bulb changes
func (l *Light) GetColourDistance() float64 {
	dist := math.Abs(float64(l.prevState.Hue) - float64(l.State.Hue))
	dist += math.Abs((float64(l.prevState.Bri) * 360.0) - (float64(l.State.Bri) * 360.0))
	dist += math.Abs((float64(l.prevState.Sat) * 360.0) - (float64(l.State.Sat) * 360.0))
	return dist
}

func (l *Light) SetStateWithTransition(transitionTime uint16) string {
	return l.setStateInternal(int32(transitionTime))
}

func (l *Light) SetState() string {
	return l.setStateInternal(0)
}

func (l *Light) setStateInternal(transitionTime int32) string {
	var buffer bytes.Buffer
	fieldsUpdated := false

	buffer.WriteString("{")
	if l.State.On != l.prevState.On {
		l.writeUpdateParam(&buffer, "on", boolToString(l.State.On), fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.Bri != l.prevState.Bri {
		l.writeUpdateParam(&buffer, "bri", fmt.Sprintf("%v", l.State.Bri), fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.Hue != l.prevState.Hue {
		l.writeUpdateParam(&buffer, "hue", fmt.Sprintf("%v", l.State.Hue), fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.Sat != l.prevState.Sat {
		l.writeUpdateParam(&buffer, "sat", fmt.Sprintf("%v", l.State.Sat), fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.X != l.prevState.X || l.State.Y != l.prevState.Y {
		l.writeUpdateParam(&buffer, "xy", "["+floatToString(l.State.X)+","+floatToString(l.State.Y)+"]", fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.Ct != l.prevState.Ct {
		l.writeUpdateParam(&buffer, "ct", fmt.Sprintf("%v", l.State.Ct), fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.Alert != l.prevState.Alert {
		l.writeUpdateParam(&buffer, "alert", stringToString(l.State.Alert), fieldsUpdated)
		fieldsUpdated = true
	}
	if l.State.Effect != l.prevState.Effect {
		l.writeUpdateParam(&buffer, "effect", stringToString(l.State.Effect), fieldsUpdated)
		fieldsUpdated = true
	}
	if transitionTime >= 0 {
		l.writeUpdateParam(&buffer, "transitiontime", fmt.Sprintf("%v", transitionTime), fieldsUpdated)
		fieldsUpdated = true
	}
	buffer.WriteString("\n}")

	l.ResetState()

	return buffer.String()
}

func boolToString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func floatToString(f float32) string {
	return fmt.Sprintf("%6.2f", f)
}

func stringToString(s string) string {
	return "\"" + s + "\""
}

func (l *Light) writeUpdateParam(buffer *bytes.Buffer, p string, val string, inclPrecedingComma bool) {
	if inclPrecedingComma {
		buffer.WriteString(",")
	}
	buffer.WriteString("\n\t\"")
	buffer.WriteString(p)
	buffer.WriteString("\": ")
	buffer.WriteString(val)
}

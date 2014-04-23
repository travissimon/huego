package huego

import (
	"bytes"
	"fmt"
	"math"
)

type State struct {
	On        bool      `json:"on"`
	Hue       uint16    `json:"hue"`
	Bri       uint8     `json:"bri"`
	Sat       uint8     `json:"sat"`
	XY        []float32 `json:"xy"`
	Ct        uint16    `json:"ct"`
	Alert     string    `json:"alert"`
	Effect    string    `json:"effect"`
	Colormode string    `json:"colormode"`
	Reachable bool      `json:"reachable"`
}

type Light struct {
	Id              string
	base            *Base
	prevState       *State
	State           *State `json:"state"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Model           string `json:"model"`
	SoftwareVersion string `json:"swversion"`
}

func (l *Light) ResetState() {
	l.prevState.On = l.State.On
	l.prevState.Hue = l.State.Hue
	l.prevState.Bri = l.State.Bri
	l.prevState.Sat = l.State.Sat
	if l.prevState.XY == nil || len(l.prevState.XY) != 2 {
		l.prevState.XY = make([]float32, 2, 2)
	}
	if l.State.XY == nil || len(l.State.XY) != 2 {
		l.State.XY = make([]float32, 2, 2)
	}
	if l.State.XY != nil && len(l.State.XY) == 2 {
		l.prevState.XY[0] = l.State.XY[0]
		l.prevState.XY[1] = l.State.XY[1]
	}
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

func (l *Light) SetStateWithTransition(transitionTime uint16) error {
	return l.setStateInternal(transitionTime)
}

func (l *Light) SetState() error {
	return l.setStateInternal(0)
}

func (l *Light) GetUpdateString(transitionTime uint16) string {
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
	if l.State.XY != nil && l.prevState.XY != nil {
		if l.State.XY[0] != l.prevState.XY[0] || l.State.XY[1] != l.prevState.XY[1] {
			l.writeUpdateParam(&buffer, "xy", "["+floatToString(l.State.XY[0])+","+floatToString(l.State.XY[1])+"]", fieldsUpdated)
			fieldsUpdated = true
		}
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

	return buffer.String()
}

func (l *Light) setStateInternal(transitionTime uint16) error {
	updateString := l.GetUpdateString(transitionTime)
	url := fmt.Sprintf("/lights/%v/state", l.Id)
	fmt.Printf("Light: %v\n", l)
	result, err := l.base.doPut(url, updateString)
	if err != nil {
		return err
	}
	resString := string(result)
	fmt.Printf("result: %v\n", resString)
	l.ResetState()
	return nil
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

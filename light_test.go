package huego

import (
	"fmt"
	"testing"
)

func Test_Update(t *testing.T) {
	l := NewLight()
	l.State.Hue = 123
	l.State.Sat = 231
	l.State.On = true
	l.State.Effect = "colorcycle"

	updateStr := l.SetStateWithTransition(123)
	fmt.Println(updateStr)

}

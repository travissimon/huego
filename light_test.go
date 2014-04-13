package huego

import (
	"fmt"
	"testing"
)

func Test_Update(t *testing.T) {
	l := NewLight()
	l.State.Hue = 123

	updateStr := l.SetState()
	fmt.Println(updateStr)
	//t.Log("update string: %s", updateStr)
}

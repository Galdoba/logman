package colorizer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gookit/color"
)

func TestColorizer(t *testing.T) {
	var colTest uint8
	for i := colTest; i <= 255; i++ {
		if i%16 == 0 {
			fmt.Println("")
		}
		col := color.S256(i, i)
		col.Print("12")
		fmt.Print(" ", i, " ")
		if i < 10 {
			fmt.Print(" ")
		}
		if i < 100 {
			fmt.Print(" ")
		}

		if i == 255 {
			fmt.Println("")
			break
		}
	}

	cl := DefaultScheme()

	for n, arg := range []interface{}{
		`//http`,
	} {
		fmt.Println("--------", n)
		s := cl.ColorizeByType(arg)
		argType := "nil"
		if arg != nil {
			argType = reflect.ValueOf(arg).Type().String()
		}

		fmt.Println(arg)
		fmt.Printf("%v' (%v %v) of type '%v'\n", s, cl.color256[fgKey(argType)], cl.color256[bgKey(argType)], argType)
	}
}

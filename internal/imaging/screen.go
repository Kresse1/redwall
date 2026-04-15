package imaging

import (
	"fmt"
	"os/exec"
	"strings"
)

type Screen struct {
	Width  int
	Height int
}

func NewScreen() (*Screen, error) {

	output, err := exec.Command("xrandr").Output()
	if err != nil {
		return nil, err
	}

	var primaryString string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, " connected primary ") {
			primaryString = line
			break
		}
	}
	if primaryString == "" {
		err = fmt.Errorf("no primary display found")
		return nil, err
	}
	fields := strings.Fields(primaryString)

	var w, h int
	for _, word := range fields {
		if strings.Contains(word, "x") && word[0] >= '0' && word[0] <= '9' {
			n, err := fmt.Sscanf(word, "%dx%d", &w, &h)
			if err != nil {
				return nil, err
			}
			if n != 2 {
				err = fmt.Errorf("no valid screen resolution found")
				return nil, err
			}
			break

		}
	}

	return &Screen{
		Width:  w,
		Height: h,
	}, nil
}

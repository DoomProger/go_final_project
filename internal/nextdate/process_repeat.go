package nextdate

import (
	"fmt"
	"strconv"
	"strings"
)

func getRepeat(repeat string) (repeatDate, error) {
	if len(repeat) == 0 {
		return repeatDate{}, fmt.Errorf("input value is empty [%s]", repeat)
	}

	repeatSettings := strings.Split(repeat, " ")

	if repeatSettings[0] == "y" && len(repeatSettings) == 1 {
		return repeatDate{
			years: 1,
		}, nil
	}

	if repeatSettings[0] == "d" && len(repeatSettings) == 2 {
		v, err := strconv.Atoi(repeatSettings[1])
		if err != nil {
			return repeatDate{}, fmt.Errorf("the number of repeating days must be define")
		}

		if v < 1 || v > 400 {
			return repeatDate{}, fmt.Errorf("the number of repeating days must be between 1 and 400, but go [%d]", v)
		}

		return repeatDate{
			days: v,
		}, nil
	}
	return repeatDate{}, fmt.Errorf("can't parse repeat values, got %v", repeatSettings)
}

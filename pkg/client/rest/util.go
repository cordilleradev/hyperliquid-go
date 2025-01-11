package rest

import (
	"fmt"
	"strconv"
)

func parseStringToFloat(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Printf("Error converting string to float64: %v", err)
		return 0.0
	}
	return value
}

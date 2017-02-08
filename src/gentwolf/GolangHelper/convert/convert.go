package convert

import (
	"fmt"
	"strconv"
)

func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func ToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func ToUint64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return i
}

func ToFloat32(str string) float32 {
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}

func ToFloat64(str string) float64 {
	i, _ := strconv.ParseFloat(str, 64)
	return i
}

func ToStr(n interface{}, args ...int) string {
	if len(args) == 0 {
		return fmt.Sprintf("%d", n)
	} else {
		format := "%." + fmt.Sprintf("%d", args[0]) + "f"
		return fmt.Sprintf(format, n)
	}
}

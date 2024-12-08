package random

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GetRandomInt(len int) int {
	return rand.Intn(9*int(math.Pow(10, float64(len-1)))) + int(math.Pow(10, float64(len-1)))
}

func GetNowAndLenRandomString(len int) string {
	return time.Now().Format("20060102") + strconv.Itoa(GetRandomInt(len))
}

package util

import (
	"strconv"
	"strings"
)

import (
	"math/rand"
)

func GenerateData(count int) []byte {
	data := make([]string, count)
	for i := 0; i < count; i++ {
		data[i] = strconv.Itoa(rand.Int() % 100)
	}
	return []byte(strings.Join(data, ","))
}

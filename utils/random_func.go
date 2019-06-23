package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Datetime func
func Datetime() string {
	now := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// HexKeyGenerator func
func HexKeyGenerator(nb int) string {
	rand.Seed(time.Now().Unix())
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	sb.Grow(nb)
	for ; nb > 0; nb-- {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes)-1)])
	}
	return sb.String()
}

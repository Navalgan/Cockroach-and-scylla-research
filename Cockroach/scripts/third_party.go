package scripts

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func NewSeed(n int64) {
	src = rand.NewSource(n)
}

func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func GetSSTFromNode(node string) []string {
	res := make([]string, 0)

	entries, err := os.ReadDir("data/" + node)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fileName := e.Name()
		if fileName[len(fileName)-3:] == "sst" {
			res = append(res, "data/"+node+"/"+fileName)
		}
	}

	return res
}

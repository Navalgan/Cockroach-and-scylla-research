package scripts

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
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

func GetProcIDs() ([]string, error) {
	res := make([]string, 0)

	lsCmd := exec.Command("ps", "aux")
	psOut, err := lsCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := lsCmd.Start(); err != nil {
		return nil, err
	}

	grepCmd := exec.Command("grep", "cockroach")
	grepCmd.Stdin = psOut

	data, err := grepCmd.Output()
	if err != nil {
		return nil, err
	}

	grepOut := string(data)

	for _, line := range strings.Split(grepOut, "\n") {
		args := strings.Split(line, " ")
		if len(args) < 5 {
			continue
		}
		res = append(res, args[4])
	}

	return res, nil
}

func GetDataFromSST(file string, grep string) (string, error) {
	out, err := exec.Command("bash", "scripts/scan.sh", file, grep).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

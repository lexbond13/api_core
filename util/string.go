package util

import (
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
)

var RegexpOnlyCarsNums = regexp.MustCompile("[^a-zA-Z0-9]+")

func RandomString(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}

// ContainsString returns true if a string is present in a iteratee.
func ContainsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func FileName(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func FileExt(fileName string) string {
	return strings.ReplaceAll(filepath.Ext(fileName), ".", "")
}

func FilterOnlyCharsNums(input string) string {
	processedString := RegexpOnlyCarsNums.ReplaceAllString(input, "")
	return processedString
}

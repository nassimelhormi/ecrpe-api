package utils

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	homePath = os.Getenv("HOME")
)

func normSubject(subjectName string) string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing 	marks
	}
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	sLower := strings.ToLower(subjectName)
	r, _, err := transform.String(t, sLower)
	if err != nil {
		logrus.Errorf("Transformer Error: %s\n", err)
		return sLower
	}
	return r
}

// MakeSessionPath func
func MakeSessionPath(data struct {
	Year string `db:"year"`
	Name string `db:"name"`
}, sessionsID int64) (string, error) {
	pathSession := fmt.Sprintf("%s/player/%s/%s/rc/session-%d", homePath, normSubject(data.Name), data.Year, sessionsID)
	if err := os.MkdirAll(pathSession, 0777); err != nil {
		return "", err
	}
	return pathSession, nil
}

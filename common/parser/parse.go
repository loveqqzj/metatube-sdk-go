package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/araddon/dateparse"
	"golang.org/x/net/html"
	dt "gorm.io/datatypes"
)

// ParseInt parses string to int regardless.
func ParseInt(s string) int {
	s = strings.TrimSpace(s)
	n, _ := strconv.ParseInt(s, 10, 64)
	return int(n)
}

// ParseTime parses a string with valid time format into time.Time.
func ParseTime(s string) time.Time {
	s = strings.TrimSpace(s)
	if ss := regexp.MustCompile(`([\s\d]+)年([\s\d]+)月([\s\d]+)日`).
		FindStringSubmatch(s); len(ss) == 4 {
		s = fmt.Sprintf("%s-%s-%s",
			strings.TrimSpace(ss[1]),
			strings.TrimSpace(ss[2]),
			strings.TrimSpace(ss[3]))
	}
	t, _ := dateparse.ParseAny(s)
	return t
}

// ParseDate parses a string with valid date format into Date.
func ParseDate(s string) dt.Date {
	return dt.Date(ParseTime(s))
}

// ParseDuration parses a string with valid duration format into time.Duration.
func ParseDuration(s string) time.Duration {
	s = ReplaceSpaceAll(s)
	s = strings.ToLower(s)
	s = strings.Replace(s, "秒", "s", 1)
	s = strings.Replace(s, "分", "m", 1)
	s = strings.Replace(s, "sec", "s", 1)
	s = strings.Replace(s, "min", "m", 1)
	if ss := regexp.MustCompile(`(?i)(\d+):(\d+):(\d+)`).FindStringSubmatch(s); len(ss) > 0 {
		s = fmt.Sprintf("%02sh%02sm%02ss", ss[1], ss[2], ss[3])
	} else if ss := regexp.MustCompile(`(?i)(\d+[mhs]?)`).FindAllStringSubmatch(s, -1); len(ss) > 0 {
		ds := make([]string, 0, 3)
		for _, d := range ss {
			ds = append(ds, d[1])
		}
		s = strings.Join(ds, "")
	}
	d, _ := time.ParseDuration(s)
	return d
}

// ParseRuntime parses a string into time.Duration and converts it to minutes as integer.
func ParseRuntime(s string) int {
	return int(ParseDuration(s).Minutes())
}

// ParseScore parses a string into float-based score.
func ParseScore(s string) float64 {
	s = strings.ReplaceAll(s, "点", "")
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return 0
	}
	s = strings.TrimSpace(fields[0])
	n, _ := strconv.ParseFloat(s, 64)
	return n
}

// ParseTexts parses all plaintext from the given *html.Node.
func ParseTexts(n *html.Node, texts *[]string) {
	if n.Type == html.TextNode {
		if text := strings.TrimSpace(n.Data); text != "" {
			*texts = append(*texts, text)
		}
	}
	for n := n.FirstChild; n != nil; n = n.NextSibling {
		ParseTexts(n, texts)
	}
}

// ReplaceSpaceAll removes all spaces in string.
func ReplaceSpaceAll(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, c := range s {
		if !unicode.IsSpace(c) {
			b.WriteRune(c)
		}
	}
	return b.String()
}

package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	linkRe = regexp.MustCompile(`(?mU)src="(http.*)"`)
)

// Slug takes in a string s and returns a user- and SEO-friendly URL.
func Slug(s string) string {
	var b strings.Builder
	var sep bool
	for _, c := range s {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'):
			b.WriteRune(unicode.ToLower(c))
			sep = true
		case c == ' ' || c == '-' || c == '_' || c == '.':
			if sep {
				b.WriteRune('-')
				sep = false
			}
		}
	}

	if b.Len() == 0 {
		return "default-slug"
	}

	return b.String()
}

func QueryUnescape(s string) string {
	s, err := url.QueryUnescape(s)
	if err != nil {
		s = err.Error()
	}
	return s
}

func ProxyResources(s, t string, id uint) string {
	sub := fmt.Sprintf(`src="/proxy?link=$1&type=%s&id=%d" loading="lazy"`, t, id)
	return linkRe.ReplaceAllString(s, sub)
}

func HumanizeNumber(i int) string {
	switch {
	case i >= 100_000:
		return format(i, 3)

	case i >= 10_000:
		return format(i, 2)

	case i >= 1_000:
		return format(i, 1)

	default:
		return strconv.Itoa(i)
	}
}

func format(i, pos int) string {
	s := fmt.Sprintf("%d", i)

	if s[pos] == '0' {
		return fmt.Sprintf("%sk", s[:pos])
	}

	return fmt.Sprintf("%s.%ck", s[:pos], s[pos])
}

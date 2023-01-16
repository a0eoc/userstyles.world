package strutils

import (
	"testing"
)

func TestSluggifyURLs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		a        string
		expected string
	}{
		{"TestProperName", "UserStyle Name", "userstyle-name"},
		{"TestMoreCharacters", "What_Even-Is  This?!", "what-even-is-this"},
		{"TestExtraCharacters", "(Dark) Something [v1.2.3]", "dark-something-v1-2-3"},
		{"TextExtraOfEverything", " Please---Get___Some   HELP!!! ", "please-get-some-help"},
		{"TestTypographicSymbols", "暗い空 Dark Mode", "dark-mode"},
		{"TestTypographicSymbolsOnly", "暗い空", "default-slug"},
	}

	for _, c := range cases {
		actual := SlugifyURL(c.a)
		if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s",
				c.desc, c.expected, actual)
		}
	}
}

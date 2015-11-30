package pcre2_test

import (
	"regexp"
	"testing"

	"github.com/lestrrat/go-pcre2"
)

type regexper interface {
	Match([]byte) bool
	MatchString(string) bool
	FindAllIndex([]byte, int) [][]int
	FindAllSubmatchIndex([]byte, int) [][]int
	FindAllStringIndex(string, int) [][]int
	FindAllStringSubmatchIndex(string, int) [][]int
	FindSubmatchIndex([]byte) []int
	FindStringSubmatchIndex(string) []int
}

func benchMatch(b *testing.B, re regexper, dos bool) {
	patterns := []string{`Hello World!`, `Hello Friend!`, `Hello 友達!`}
	for _, pat := range patterns {
		var rv bool
		if dos {
			rv = re.MatchString(pat)
		} else {
			rv = re.Match([]byte(pat))
		}
		if !rv {
			b.Errorf("Expected to match, failed")
			return
		}
	}

	patterns = []string{`Goodbye World!`, `Hell no!`, `HelloWorld!`}
	for _, pat := range patterns {
		var rv bool
		if dos {
			rv = re.MatchString(pat)
		} else {
			rv = re.Match([]byte(pat))
		}

		if rv {
			b.Errorf("Expected to NOT match, matched")
			return
		}
	}
}

func benchFindAllIndex(b *testing.B, re regexper, dos bool) {
	patterns := []string{`Alice:35 Bob:42 Charlie:21`, `桃:三年 栗:三年 柿:八年`, `vini:came vidi:saw vici:won`}
	for _, pat := range patterns {
		var matches [][]int
		if dos {
			matches = re.FindAllStringIndex(pat, -1)
		} else {
			matches = re.FindAllIndex([]byte(pat), -1)
		}

		if len(matches) != 3 {
			b.Errorf("Expected to match '%s' against '%#v', got %d", pat, re, len(matches))
			b.Logf("%#v", matches)
			return
		}
	}
}

func benchFindSubmatchIndex(b *testing.B, re regexper, dos bool) {
	patterns := []string{`Alice:35 Bob:42 Charlie:21`, `桃:三年 栗:三年 柿:八年`, `vini:came vidi:saw vici:won`}
	for _, pat := range patterns {
		var matches []int
		if dos {
			matches = re.FindStringSubmatchIndex(pat)
		} else {
			matches = re.FindSubmatchIndex([]byte(pat))
		}

		if len(matches) != 6 {
			b.Errorf("Expected to match '%s' against '%#v', got %d", pat, re, len(matches))
			b.Logf("%#v", matches)
			return
		}
	}
}

func benchFindAllSubmatchIndex(b *testing.B, re regexper, dos bool) {
	patterns := []string{`Alice:35 Bob:42 Charlie:21`, `桃:三年 栗:三年 柿:八年`, `vini:came vidi:saw vici:won`}
	for _, pat := range patterns {
		var matches [][]int
		if dos {
			matches = re.FindAllStringSubmatchIndex(pat, -1)
		} else {
			matches = re.FindAllSubmatchIndex([]byte(pat), -1)
		}

		if len(matches) != 3 {
			b.Errorf("Expected to match '%s' against '%#v', got %d", pat, re, len(matches))
			b.Logf("%#v", matches)
			return
		}
	}
}

func makeBenchFunc(b *testing.B, which bool, dos bool, pattern string, f func(*testing.B, regexper, bool)) func() {
	// Forcing a function call so that we have chance to
	// run garbage collection for each iteration
	return func() {
		var re regexper
		var err error
		if which { // true == pcre2
			re, err = pcre2.Compile(pattern)
		} else {
			re, err = regexp.Compile(pattern)
		}
		if err != nil {
			b.Errorf("compile failed: %s", err)
			return
		}
		f(b, re, dos)
	}
}

const (
	UseGoRegexp    = false
	UsePCRE2Regexp = true
	UseBytes       = false
	UseString      = true
)

// Match, MatchString
const RegexpMatchRegex = `^Hello (.+)!$`

func BenchmarkGoRegexpMatch(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseBytes, RegexpMatchRegex, benchMatch)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2RegexpMatch(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseBytes, RegexpMatchRegex, benchMatch)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkGoRegexpMatchString(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseString, RegexpMatchRegex, benchMatch)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2RegexpMatchString(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseString, RegexpMatchRegex, benchMatch)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

// FindAllIndex, FindAllStringIndex
const FindAllIndexRegex = `(\S+):(\S+)`

func BenchmarkGoFindAllIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseBytes, FindAllIndexRegex, benchFindAllIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindAllIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseBytes, FindAllIndexRegex, benchFindAllIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkGoFindAllStringIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseString, FindAllIndexRegex, benchFindAllIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindAllStringIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseString, FindAllIndexRegex, benchFindAllIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

// FindSubmatchIndex, FindStringSubmatchIndex
const FindSubmatchIndexRegex = `(\S+):(\S+)`

func BenchmarkGoFindSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseBytes, FindSubmatchIndexRegex, benchFindSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseBytes, FindSubmatchIndexRegex, benchFindSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkGoFindStringSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseString, FindSubmatchIndexRegex, benchFindSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindStringSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseString, FindSubmatchIndexRegex, benchFindSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

// FindAllSubmatchIndex, FindAllStringSubmatchIndex
const FindAllSubmatchIndexRegex = `(\S+):(\S+)`

func BenchmarkGoFindAllSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseBytes, FindAllSubmatchIndexRegex, benchFindAllSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindAllSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseBytes, FindAllSubmatchIndexRegex, benchFindAllSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkGoFindAllStringSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UseGoRegexp, UseString, FindAllSubmatchIndexRegex, benchFindAllSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindAllStringSubmatchIndex(b *testing.B) {
	benchf := makeBenchFunc(b, UsePCRE2Regexp, UseString, FindAllSubmatchIndexRegex, benchFindAllSubmatchIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

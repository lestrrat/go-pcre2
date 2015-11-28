package pcre2_test

import (
	"regexp"
	"testing"

	"github.com/lestrrat/go-pcre2"
)

type regexper interface {
	MatchString(string) bool
	FindAllIndex([]byte, int) [][]int
}

func benchMatchString(b *testing.B, re regexper) {
	patterns := []string{`Hello World!`, `Hello Friend!`, `Hello 友達!`}
	for _, pat := range patterns {
		if !re.MatchString(pat) {
			b.Errorf("Expected to match, failed")
			return
		}
	}

	patterns = []string{`Goodbye World!`, `Hell no!`, `HelloWorld!`}
	for _, pat := range patterns {
		if re.MatchString(pat) {
			b.Errorf("Expected to NOT match, matched")
			return
		}
	}
}

func benchFindAllIndex(b *testing.B, re regexper) {
	patterns := []string{`Alice Bob Charlie`, `桃 栗 柿`, `vini vidi vici`}
	for _, pat := range patterns {
		matches := re.FindAllIndex([]byte(pat), -1)
		if len(matches) != 3 {
			b.Errorf("Expected to match '%s' against '%#v', got %d", pat, re, len(matches))
			b.Logf("%#v", matches)
			return
		}
	}
}

func makeGoBenchFunc(b *testing.B, pattern string, f func(*testing.B, regexper)) func() {
	// Forcing a function call so that we have chance to
	// run garbage collection for each iteration
	return func() {
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Errorf("compile failed: %s", err)
			return
		}
		f(b, re)
	}
}

func makePCRE2BenchFunc(b *testing.B, pattern string, f func(*testing.B, regexper)) func() {
	// Forcing a function call so that we have chance to
	// run garbage collection for each iteration
	return func() {
		re, err := pcre2.Compile(pattern)
		if err != nil {
			b.Errorf("compile failed: %s", err)
			return
		}
		defer re.Free()
		f(b, re)
	}
}

const RegexpMatchRegex = `^Hello (.+)!$`
func BenchmarkGoRegexpMatch(b *testing.B) {
	benchf := makeGoBenchFunc(b, RegexpMatchRegex, benchMatchString)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2RegexpMatch(b *testing.B) {
	benchf := makePCRE2BenchFunc(b, RegexpMatchRegex, benchMatchString)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

const FindAllIndexRegex = `(\S+)`
func BenchmarkGoFindAllIndex(b *testing.B) {
	benchf := makeGoBenchFunc(b, FindAllIndexRegex, benchFindAllIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}

func BenchmarkPCRE2FindAllIndex(b *testing.B) {
	benchf := makePCRE2BenchFunc(b, FindAllIndexRegex, benchFindAllIndex)
	for i := 0; i < b.N; i++ {
		benchf()
	}
}


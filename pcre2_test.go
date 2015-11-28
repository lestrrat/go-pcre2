package pcre2_test

import (
	"testing"

	"github.com/lestrrat/go-pcre2"
	"github.com/stretchr/testify/assert"
)

func TestBadPattern(t *testing.T) {
	re, err := pcre2.Compile(`^Hello [World!$`)
	if !assert.Error(t, err, "Compile works") {
		return
	}
	defer re.Free()
}

func TestBasic(t *testing.T) {
	re, err := pcre2.Compile(`^Hello (.+)!$`)
	if !assert.NoError(t, err, "Compile works") {
		return
	}
	defer re.Free()

	patterns := []string{`Hello World!`, `Hello Friend!`, `Hello 友達!`}
	for _, pat := range patterns {
		t.Logf("Matching against []byte '%s' (expect MATCH)", pat)
		if !assert.True(t, re.Match([]byte(pat)), "Match succeeds for %s", pat) {
			return
		}

		t.Logf("Matching against string '%s' (expect MATCH)", pat)
		if !assert.True(t, re.MatchString(pat), "MatchString succeeds for %s", pat) {
			return
		}
	}

	patterns = []string{`Goodbye World!`, `Hell no!`, `HelloWorld!`}
	for _, pat := range patterns {
		t.Logf("Matching against []byte '%s' (expect FAIL)", pat)
		if !assert.False(t, re.Match([]byte(pat)), "Match fails for %s", pat) {
			return
		}

		t.Logf("Matching against string '%s' (expect FAIL)", pat)
		if !assert.False(t, re.MatchString(pat), "MatchString fails for %s", pat) {
			return
		}
	}
}

func TestFindAllSubmatchIndex(t *testing.T) {
	re, err := pcre2.Compile(`^(\S+) (\S+) (\S+)$`)
	if !assert.NoError(t, err, "Compile works") {
		return
	}
	defer re.Free()

	data := map[string][][]int{
		`Alice Bob Charlie`: [][]int{[]int{0, 17}, []int{0, 5}, []int{6, 9}, []int{10, 17}},
		`桃 栗 柿`:             [][]int{[]int{0, 5}, []int{0, 1}, []int{2, 3}, []int{4, 5}},
		`vini vidi vici`:    [][]int{[]int{0, 14}, []int{0, 4}, []int{5, 9}, []int{10, 14}},
	}
	for subject, expected := range data {
		ret := re.FindAllSubmatchIndex([]byte(subject), 0)
		if !assert.NotEmpty(t, ret, "Match should succeed") {
			return
		}

		if !assert.Equal(t, expected, ret, "indices should match") {
			return
		}
	}
}
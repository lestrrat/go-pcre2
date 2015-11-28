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
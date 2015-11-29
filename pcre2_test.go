package pcre2_test

import (
	"regexp"
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

func TestFindAllIndex(t *testing.T) {
	pattern := `(\S+):(\S+)`
	gore, err := regexp.Compile(pattern)
	if !assert.NoError(t, err, "Compile works (Go)") {
		return
	}

	re, err := pcre2.Compile(pattern)
	if !assert.NoError(t, err, "Compile works (pcre2)") {
		return
	}
	defer re.Free()

	data := []string{`Alice:35 Bob:42 Charlie:21`, `桃:三年 栗:三年 柿:八年`, `vini:came vidi:saw vici:won`}
	for _, subject := range data {
		t.Logf("FindAllIndex against '%s'", subject)
		expected := gore.FindAllIndex([]byte(subject), -1)
		ret := re.FindAllIndex([]byte(subject), -1)
		if !assert.NotEmpty(t, ret, "Match should succeed") {
			return
		}

		if !assert.Equal(t, expected, ret, "indices should match") {
			return
		}
	}
}

func TestFindAllSubmatchIndex(t *testing.T) {
	pattern := `(\S+):(\S+)`
	gore, err := regexp.Compile(pattern)
	if !assert.NoError(t, err, "Compile works (Go)") {
		return
	}

	re, err := pcre2.Compile(pattern)
	if !assert.NoError(t, err, "Compile works (pcre2)") {
		return
	}
	defer re.Free()

	data := []string{`Alice:35 Bob:42 Charlie:21`, `桃:三年 栗:三年 柿:八年`, `vini:came vidi:saw vici:won`}
	for _, subject := range data {
		t.Logf("FindAllSubmatchIndex against '%s'", subject)
		expected := gore.FindAllSubmatchIndex([]byte(subject), -1)
		ret := re.FindAllSubmatchIndex([]byte(subject), -1)
		if !assert.NotEmpty(t, ret, "Match should succeed") {
			return
		}

		if !assert.Equal(t, expected, ret, "indices should match") {
			return
		}
	}
}

func TestFindAllStringSubmatchIndex(t *testing.T) {
	pattern := `(\S+):(\S+)`
	gore, err := regexp.Compile(pattern)
	if !assert.NoError(t, err, "Compile works (Go)") {
		return
	}

	re, err := pcre2.Compile(pattern)
	if !assert.NoError(t, err, "Compile works (pcre2)") {
		return
	}
	defer re.Free()

	data := []string{`Alice:35 Bob:42 Charlie:21`, `桃:三年 栗:三年 柿:八年`, `vini:came vidi:saw vici:won`}
	for _, subject := range data {
		t.Logf("FindAllStringSubmatchIndex against '%s'", subject)
		expected := gore.FindAllStringSubmatchIndex(subject, -1)
		ret := re.FindAllStringSubmatchIndex(subject, -1)
		if !assert.NotEmpty(t, ret, "Match should succeed") {
			return
		}

		if !assert.Equal(t, expected, ret, "indices should match") {
			return
		}
	}
}






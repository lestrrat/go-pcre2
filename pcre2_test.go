package pcre2_test

import (
	"regexp"
	"testing"

	"github.com/lestrrat/go-pcre2"
	"github.com/stretchr/testify/assert"
)

func TestBadPattern(t *testing.T) {
	re, err := pcre2.Compile(`^Hello [World!$`)
	t.Logf("%s", err)
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

func TestFind(t *testing.T) {
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
	for _, doString := range []bool{true, false} {
		var methodName string
		if doString {
			methodName = "FindString"
		} else {
			methodName = "Find"
		}

		var expected interface{}
		var ret interface{}

		for _, subject := range data {
			t.Logf(`%s("%s")`, methodName, subject)
			if doString {
				expected = gore.FindString(subject)
				ret = re.FindString(subject)
			} else {
				expected = gore.Find([]byte(subject))
				ret = re.Find([]byte(subject))
			}

			if !assert.Equal(t, expected, ret, "returned byte sequence should match") {
				return
			}
		}
	}
}

func TestFindIndex(t *testing.T) {
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
	for _, doString := range []bool{true, false} {
		var methodName string
		if doString {
			methodName = "FindStringIndex"
		} else {
			methodName = "FindIndex"
		}

		var expected interface{}
		var ret interface{}

		for _, subject := range data {
			t.Logf(`%s("%s")`, methodName, subject)
			if doString {
				expected = gore.FindStringIndex(subject)
				ret = re.FindStringIndex(subject)
			} else {
				expected = gore.FindIndex([]byte(subject))
				ret = re.FindIndex([]byte(subject))
			}

			if !assert.Equal(t, expected, ret, "returned byte sequence should match") {
				return
			}
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
	for _, doString := range []bool{true, false} {
		var methodName string
		if doString {
			methodName = "FindAllStringIndex"
		} else {
			methodName = "FindAllIndex"
		}

		for n := -1; n < 4; n++ {
			var expected [][]int
			var ret [][]int

			for _, subject := range data {
				t.Logf(`%s("%s", %d)`, methodName, subject, n)
				if doString {
					expected = gore.FindAllStringIndex(subject, n)
					ret = re.FindAllStringIndex(subject, n)
				} else {
					expected = gore.FindAllIndex([]byte(subject), n)
					ret = re.FindAllIndex([]byte(subject), n)
				}

				if !assert.Equal(t, expected, ret, "indices should match") {
					return
				}
			}
		}
	}
}

func TestFindAll(t *testing.T) {
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
	for _, doString := range []bool{true, false} {
		var methodName string
		if doString {
			methodName = "FindAllString"
		} else {
			methodName = "FindAll"
		}

		for n := -1; n < 4; n++ {
			var expected interface{}
			var ret interface{}

			for _, subject := range data {
				t.Logf(`%s("%s", %d)`, methodName, subject, n)
				if doString {
					expected = gore.FindAllString(subject, n)
					ret = re.FindAllString(subject, n)
				} else {
					expected = gore.FindAll([]byte(subject), n)
					ret = re.FindAll([]byte(subject), n)
				}

				if !assert.Equal(t, expected, ret, "indices should match") {
					return
				}
			}
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
	for _, doString := range []bool{true, false} {
		var methodName string
		if doString {
			methodName = "FindAllStringSubmatchIndex"
		} else {
			methodName = "FindAllSubmatcnIndex"
		}

		for n := -1; n < 4; n++ {
			var expected [][]int
			var ret [][]int

			for _, subject := range data {
				t.Logf(`%s("%s", %d)`, methodName, subject, n)
				if doString {
					expected = gore.FindAllStringSubmatchIndex(subject, n)
					ret = re.FindAllStringSubmatchIndex(subject, n)
				} else {
					expected = gore.FindAllSubmatchIndex([]byte(subject), n)
					ret = re.FindAllSubmatchIndex([]byte(subject), n)
				}

				if !assert.Equal(t, expected, ret, "indices should match") {
					return
				}
			}
		}
	}
}

func TestFindAllSubmatch(t *testing.T) {
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
	for _, doString := range []bool{true, false} {
		var methodName string
		if doString {
			methodName = "FindAllStringSubmatch"
		} else {
			methodName = "FindAllSubmatcn"
		}

		for n := -1; n < 4; n++ {
			var expected interface{}
			var ret interface{}

			for _, subject := range data {
				t.Logf(`%s("%s", %d)`, methodName, subject, n)
				if doString {
					expected = gore.FindAllStringSubmatch(subject, n)
					ret = re.FindAllStringSubmatch(subject, n)
				} else {
					expected = gore.FindAllSubmatch([]byte(subject), n)
					ret = re.FindAllSubmatch([]byte(subject), n)
				}

				if !assert.Equal(t, expected, ret, "indices should match") {
					return
				}
			}
		}
	}
}


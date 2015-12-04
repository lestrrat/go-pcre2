/*
Package pcre2 is a wrapper around PCRE2 C library. This library aims to
provide compatible API as that of regexp package from Go stdlib.

Note that while PCRE2 provides support for 8, 16, and 32 bit inputs,
this library assumes UTF-8 (32bit) input. Therefore if you use anything
other than UTF-8, matches will not succeed.
*/
package pcre2

/*
#define PCRE2_CODE_UNIT_WIDTH 32
#cgo pkg-config: libpcre2-32
#include <stdio.h>
#include <stdlib.h>
#include <pcre2.h>

#define MY_PCRE2_ERROR_MESSAGE_BUF_LEN 256
static
void *
MY_pcre2_get_error_message(int errnum) {
	PCRE2_UCHAR *buf = (PCRE2_UCHAR *) malloc(sizeof(PCRE2_UCHAR) * MY_PCRE2_ERROR_MESSAGE_BUF_LEN);
  pcre2_get_error_message(errnum, buf, MY_PCRE2_ERROR_MESSAGE_BUF_LEN);
	return buf;
}

*/
import "C"
import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

// Error returns the string representation of the error.
func (e ErrCompile) Error() string {
	return fmt.Sprintf("PCRE2 compilation failed at offset %d: %s", e.offset, e.message)
}

func strToRuneArray(s string) ([]rune, []int, error) {
	rs := []rune{}
	ls := []int{} // length of each rune
	for len(s) > 0 {
		r, n := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return nil, nil, ErrInvalidUTF8String
		}
		s = s[n:]
		rs = append(rs, r)
		ls = append(ls, n)
	}
	return rs, ls, nil
}

func bytesToRuneArray(b []byte) ([]rune, []int, error) {
	rs := []rune{} // actual runes
	ls := []int{}  // length of each rune
	for len(b) > 0 {
		r, n := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			return nil, nil, ErrInvalidUTF8String
		}
		b = b[n:]
		rs = append(rs, r)
		ls = append(ls, n)
	}
	return rs, ls, nil
}

// Compile takes the input string and creates a compiled Regexp object.
// Regexp objects created by Compile must be released by calling Free
func Compile(pattern string) (*Regexp, error) {
	patc, _, err := strToRuneArray(pattern)
	if err != nil {
		return nil, err
	}

	var errnum C.int
	var erroff C.PCRE2_SIZE
	re := C.pcre2_compile(
		(C.PCRE2_SPTR)(unsafe.Pointer(&patc[0])),
		C.size_t(len(patc)),
		0,
		&errnum,
		&erroff,
		nil,
	)
	if re == nil {
		rawbytes := C.MY_pcre2_get_error_message(errnum)
		msg := C.GoBytes(rawbytes, 32/8*256)
		defer C.free(unsafe.Pointer(rawbytes))

		return nil, ErrCompile{
			pattern: pattern,
			offset:  int(erroff),
			message: string(msg),
		}
	}
	return &Regexp{
		pattern: pattern,
		ptr:     uintptr(unsafe.Pointer(re)),
	}, nil
}

// MustCompile is like Compile but panics if the expression cannot be
// parsed.
func MustCompile(pattern string) *Regexp {
	r, err := Compile(pattern)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Regexp) validRegexpPtr() (*C.pcre2_code, error) {
	if r == nil {
		return nil, ErrInvalidRegexp
	}

	if rptr := r.ptr; rptr != 0 {
		return (*C.pcre2_code)(unsafe.Pointer(rptr)), nil
	}
	return nil, ErrInvalidRegexp
}

// Free releases the underlying C resources
func (r *Regexp) Free() error {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return err
	}
	C.pcre2_code_free(rptr)
	r.ptr = 0
	return nil
}

// String returns the source text used to compile the regular expression.
func (r Regexp) String() string {
	return r.pattern
}

func (r *Regexp) Match(b []byte) bool {
	rs, _, err := bytesToRuneArray(b)
	if err != nil {
		return false
	}
	return r.matchRuneArray(rs, 0, 0, nil) >= 0
}

func (r *Regexp) MatchString(s string) bool {
	rs, _, err := strToRuneArray(s)
	if err != nil {
		return false
	}
	return r.matchRuneArray(rs, 0, 0, nil) >= 0
}

func (r *Regexp) matchRuneArray(rs []rune, offset int, options int, matchData *C.pcre2_match_data) int {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return -1
	}

	if matchData == nil {
		matchData = C.pcre2_match_data_create_from_pattern(rptr, nil)
		defer C.pcre2_match_data_free(matchData)
	}

	rc := C.pcre2_match(
		rptr,
		(C.PCRE2_SPTR)(unsafe.Pointer(&rs[0])),
		C.size_t(len(rs)),
		(C.PCRE2_SIZE)(offset),
		(C.uint32_t)(options),
		matchData,
		nil,
	)

	return int(rc)
}

func pcre2GetOvectorPointer(matchData *C.pcre2_match_data, howmany int) []C.size_t {
	ovector := C.pcre2_get_ovector_pointer(matchData)
	// Note that by doing this SliceHeader maigc, we allow Go
	// slice syntax but Go doesn't own the underlying pointer.
	// We need to free it. In this case, it means the caller
	// must remember to free matchData
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ovector)),
		Len:  howmany * 2,
		Cap:  howmany * 2,
	}
	return *(*[]C.size_t)(unsafe.Pointer(&hdr))
}

func (r *Regexp) HasOption(opt int) bool {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return false
	}

	var i C.uint32_t
	C.pcre2_pattern_info(rptr, C.PCRE2_INFO_ALLOPTIONS, unsafe.Pointer(&i))
	return (uint32(i) & uint32(opt)) != 0
}

func (r *Regexp) isCRLFValid() bool {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return false
	}

	var i C.uint32_t
	C.pcre2_pattern_info(rptr, C.PCRE2_INFO_NEWLINE, unsafe.Pointer(&i))
	switch i {
	case C.PCRE2_NEWLINE_ANY, C.PCRE2_NEWLINE_CRLF, C.PCRE2_NEWLINE_ANYCRLF:
		return true
	}

	return false
}

func (r *Regexp) FindIndex(b []byte) []int {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
		return nil
	}

	is := r.findAllIndex(rs, ls, 1)
	if len(is) != 1 {
		return nil
	}
	return is[0]
}

func (r *Regexp) Find(b []byte) []byte {
	is := r.FindIndex(b)
	if is == nil {
		return nil
	}
	return b[is[0]:is[1]]
}

func (r *Regexp) FindStringIndex(s string) []int {
	rs, ls, err := strToRuneArray(s)
	if err != nil {
		return nil
	}

	is := r.findAllIndex(rs, ls, 1)
	if len(is) != 1 {
		return nil
	}
	return is[0]
}

func (r *Regexp) FindSubmatch(b []byte) [][]byte {
	matches := r.FindSubmatchIndex(b)
	if matches == nil {
		return nil
	}

	ret := make([][]byte, 0, len(matches)/2)
	for i := 0; i < len(matches)/2; i++ {
		ret = append(ret, b[matches[2*i]:matches[2*i+1]])
	}
	return ret
}

func (r *Regexp) FindSubmatchIndex(b []byte) []int {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
		return nil
	}
	return r.findSubmatchIndex(rs, ls)
}

func (r *Regexp) FindStringSubmatchIndex(s string) []int {
	rs, ls, err := strToRuneArray(s)
	if err != nil {
		return nil
	}
	return r.findSubmatchIndex(rs, ls)
}

func (r *Regexp) findSubmatchIndex(rs []rune, ls []int) []int {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return nil
	}

	matchData := C.pcre2_match_data_create_from_pattern(rptr, nil)
	defer C.pcre2_match_data_free(matchData)

	out := []int(nil)
	options := 0

	count := r.matchRuneArray(rs, 0, options, matchData)
	if count <= 0 {
		return nil
	}

	ovector := pcre2GetOvectorPointer(matchData, count)
	for i := 0; i < count; i++ {
		ovec0 := int(ovector[2*i])
		b1 := 0
		for x := 0; x < ovec0; x++ {
			b1 += ls[x]
		}

		ovec1 := int(ovector[2*i+1])
		b2 := b1
		for x := ovec0; x < ovec1; x++ {
			b2 += ls[x]
		}
		out = append(out, []int{b1, b2}...)
	}

	return out
}

func (r *Regexp) FindStringSubmatch(s string) []string {
	matches := r.FindStringSubmatchIndex(s)
	if matches == nil {
		return nil
	}

	ret := make([]string, 0, len(matches))
	for i := 0; i < len(matches)/2; i++ {
		ret = append(ret, s[matches[2*i]:matches[2*i+1]])
	}
	return ret
}

func (r *Regexp) FindString(s string) string {
	is := r.FindStringIndex(s)
	if is == nil {
		return ""
	}
	return s[is[0]:is[1]]
}

func (r *Regexp) FindAll(b []byte, n int) [][]byte {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
		return nil
	}
	ret := [][]byte(nil)
	for _, is := range r.findAllIndex(rs, ls, n) {
		ret = append(ret, b[is[0]:is[1]])
	}
	return ret
}

func (r *Regexp) FindAllString(s string, n int) []string {
	if n == 0 {
		return nil
	}

	rs, ls, err := strToRuneArray(s)
	if err != nil {
		return nil
	}
	ret := []string{}
	for _, is := range r.findAllIndex(rs, ls, n) {
		ret = append(ret, s[is[0]:is[1]])
		if n > 0 && len(ret) >= n {
			break
		}
	}
	return ret
}

func (r *Regexp) findAllIndex(rs []rune, ls []int, n int) [][]int {
	if n == 0 {
		return nil
	}

	rptr, err := r.validRegexpPtr()
	if err != nil {
		return nil
	}

	matchData := C.pcre2_match_data_create_from_pattern(rptr, nil)
	defer C.pcre2_match_data_free(matchData)

	out := [][]int(nil)
	offset := 0
	options := 0
	for len(rs) > 0 {
		count := r.matchRuneArray(rs, 0, options, matchData)
		if count <= 0 {
			break
		}

		ovector := pcre2GetOvectorPointer(matchData, count)
		ovec0 := int(ovector[0])
		b1 := 0
		for x := 0; x < ovec0; x++ {
			b1 += ls[x]
		}
		b2 := b1
		for x := ovec0; x < int(ovector[1]); x++ {
			b2 += ls[x]
		}
		out = append(out, []int{offset + b1, offset + b2})
		units := int(ovector[1])
		for x := 0; x < units; x++ {
			offset += ls[x]
		}

		rs = rs[units:]
		ls = ls[units:]

		if n > 0 && len(out) >= n {
			break
		}
	}

	return out
}

func (r *Regexp) FindAllIndex(b []byte, n int) [][]int {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
		return nil
	}
	return r.findAllIndex(rs, ls, n)
}

func (r *Regexp) FindAllStringIndex(s string, n int) [][]int {
	rs, ls, err := strToRuneArray(s)
	if err != nil {
		return nil
	}
	return r.findAllIndex(rs, ls, n)
}

func (r *Regexp) findAllSubmatchIndex(rs []rune, ls []int, n int) [][]int {
	if n == 0 {
		return nil
	}

	rptr, err := r.validRegexpPtr()
	if err != nil {
		return nil
	}

	matchData := C.pcre2_match_data_create_from_pattern(rptr, nil)
	defer C.pcre2_match_data_free(matchData)

	out := [][]int(nil)
	offset := 0
	options := 0
	for len(rs) > 0 {
		count := r.matchRuneArray(rs, 0, options, matchData)
		if count <= 0 {
			break
		}

		ovector := pcre2GetOvectorPointer(matchData, count)
		curmatch := make([]int, 0, count)
		for i := 0; i < count; i++ {
			ovec2i := int(ovector[2*i])

			b1 := 0
			for x := 0; x < ovec2i; x++ {
				b1 += ls[x]
			}
			b2 := b1
			for x := ovec2i; x < int(ovector[2*i+1]); x++ {
				b2 += ls[x]
			}
			curmatch = append(curmatch, offset+b1, offset+b2)
		}
		out = append(out, curmatch)

		units := int(ovector[1])
		for x := 0; x < units; x++ {
			offset += ls[x]
		}

		rs = rs[units:]
		ls = ls[units:]

		if n > 0 && len(out) >= n {
			break
		}
	}

	return out
}

func (r *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
		return nil
	}

	all := r.findAllSubmatchIndex(rs, ls, n)
	if all == nil {
		return nil
	}

	ret := make([][][]byte, 0, len(all))
	for _, is := range all {
		l := len(is) / 2
		cur := make([][]byte, 0, l)
		for i := 0; i < l; i++ {
			cur = append(cur, b[is[2*i]:is[2*i+1]])
		}

		ret = append(ret, cur)
		if n > 0 && len(ret) >= n {
			break
		}
	}
	return ret
}

func (r *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	rs, ls, err := strToRuneArray(s)
	if err != nil {
		return nil
	}

	all := r.findAllSubmatchIndex(rs, ls, n)
	if all == nil {
		return nil
	}

	ret := make([][]string, 0, len(all))
	for _, is := range all {
		l := len(is) / 2
		cur := make([]string, 0, l)
		for i := 0; i < l; i++ {
			cur = append(cur, s[is[2*i]:is[2*i+1]])
		}
		ret = append(ret, cur)
		if n > 0 && len(ret) >= n {
			break
		}
	}
	return ret
}

func (r *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
		return nil
	}
	return r.findAllSubmatchIndex(rs, ls, n)
}

func (r *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int {
	rs, ls, err := strToRuneArray(s)
	if err != nil {
		return nil
	}
	return r.findAllSubmatchIndex(rs, ls, n)
}

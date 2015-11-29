package pcre2

/*
#define PCRE2_CODE_UNIT_WIDTH 32
#cgo pkg-config: libpcre2-32
#include <stdlib.h>
#include <pcre2.h>

#define MY_PCRE2_ERROR_MESSAGE_BUF_LEN 256
static
char *
MY_pcre2_get_error_message(int errno) {
	PCRE2_UCHAR *buf = (PCRE2_UCHAR *) malloc(sizeof(PCRE2_UCHAR) * MY_PCRE2_ERROR_MESSAGE_BUF_LEN);
  pcre2_get_error_message(errno, buf, MY_PCRE2_ERROR_MESSAGE_BUF_LEN);
	return (char *) buf;
}

*/
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

type Regexp struct {
	ptr *C.pcre2_code
}

var (
	ErrInvalidRegexp     = errors.New("invalid regexp")
	ErrInvalidUTF8String = errors.New("invalid utf8 string")
)

type ErrCompile struct {
	message string
	offset  int
	pattern string
}

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
		// note: malloc'ed, but Go should be able to free it
		msg := C.MY_pcre2_get_error_message(errnum)
		return nil, ErrCompile{
			pattern: pattern,
			offset:  int(erroff),
			message: C.GoString(msg),
		}
	}
	return &Regexp{ptr: re}, nil
}

func (r *Regexp) validRegexpPtr() (*C.pcre2_code, error) {
	if r == nil {
		return nil, ErrInvalidRegexp
	}

	rptr := r.ptr
	if r.ptr == nil {
		return nil, ErrInvalidRegexp
	}
	return rptr, nil
}

func (r *Regexp) Free() error {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return err
	}
	C.pcre2_code_free(rptr)
	return nil
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

func (r *Regexp) FindAllIndex(b []byte, n int) [][]int {
	rs, ls, err := bytesToRuneArray(b)
	if err != nil {
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
		//		out = append(out, curmatch)

		rs = rs[units:]
		ls = ls[units:]
	}

	return out
}

func (r *Regexp) findAllSubmatchIndex(rs []rune, ls []int, n int) [][]int {
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
	}

	return out
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


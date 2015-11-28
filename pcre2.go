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

func strToRuneArray(s string) ([]rune, error) {
	rs := []rune{}
	for len(s) > 0 {
		r, n := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return nil, ErrInvalidUTF8String
		}
		s = s[n:]
		rs = append(rs, r)
	}
	return rs, nil
}

func Compile(pattern string) (*Regexp, error) {
	patc, err := strToRuneArray(pattern)
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

func (r *Regexp) MatchString(s string) bool {
	rptr, err := r.validRegexpPtr()
	if err != nil {
		return false
	}

	sc, err := strToRuneArray(s)
	if err != nil {
		return false
	}

	match_data := C.pcre2_match_data_create_from_pattern(rptr, nil)
	defer C.pcre2_match_data_free(match_data)

	rc := C.pcre2_match(
		rptr,
		(C.PCRE2_SPTR)(unsafe.Pointer(&sc[0])),
		C.size_t(len(sc)),
		0,
		0,
		match_data,
		nil,
	)

	return int(rc) >= 0
}

//re = pcre2_compile(
//  pattern,               /* the pattern */
//  PCRE2_ZERO_TERMINATED, /* indicates pattern is zero-terminated */
//  0,                     /* default options */
//  &errornumber,          /* for error number */
//  &erroroffset,          /* for error offset */
//  NULL);                 /* use default compile context */
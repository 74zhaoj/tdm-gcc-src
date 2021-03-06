// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

// The Error interface identifies a run time error.
type Error interface {
	error

	// RuntimeError is a no-op function but
	// serves to distinguish types that are run time
	// errors from ordinary errors: a type is a
	// run time error if it has a RuntimeError method.
	RuntimeError()
}

// A TypeAssertionError explains a failed type assertion.
type TypeAssertionError struct {
	_interface    *_type
	concrete      *_type
	asserted      *_type
	missingMethod string // one method needed by Interface, missing from Concrete
}

func (*TypeAssertionError) RuntimeError() {}

func (e *TypeAssertionError) Error() string {
	inter := "interface"
	if e._interface != nil {
		inter = e._interface.string()
	}
	as := e.asserted.string()
	if e.concrete == nil {
		return "interface conversion: " + inter + " is nil, not " + as
	}
	cs := e.concrete.string()
	if e.missingMethod == "" {
		msg := "interface conversion: " + inter + " is " + cs + ", not " + as
		if cs == as {
			// provide slightly clearer error message
			if e.concrete.pkgpath() != e.asserted.pkgpath() {
				msg += " (types from different packages)"
			} else {
				msg += " (types from different scopes)"
			}
		}
		return msg
	}
	return "interface conversion: " + cs + " is not " + as +
		": missing method " + e.missingMethod
}

// Remove quoted strings from gccgo reflection strings.
func unquote(s string) string {
	ls := len(s)
	var i int
	for i = 0; i < ls; i++ {
		if s[i] == '\t' {
			break
		}
	}
	if i == ls {
		return s
	}
	var q bool
	r := make([]byte, len(s))
	j := 0
	for i = 0; i < ls; i++ {
		if s[i] == '\t' {
			q = !q
		} else if !q {
			r[j] = s[i]
			j++
		}
	}
	return string(r[:j])
}

// An errorString represents a runtime error described by a single string.
type errorString string

func (e errorString) RuntimeError() {}

func (e errorString) Error() string {
	return "runtime error: " + string(e)
}

// An errorCString represents a runtime error described by a single C string.
// Not "type errorCString uintptr" because of http://golang.org/issue/7084.
type errorCString struct{ cstr uintptr }

func (e errorCString) RuntimeError() {}

func (e errorCString) Error() string {
	return "runtime error: " + gostringnocopy((*byte)(unsafe.Pointer(e.cstr)))
}

// For calling from C.
func NewErrorCString(s uintptr, ret *interface{}) {
	*ret = errorCString{s}
}

// plainError represents a runtime error described a string without
// the prefix "runtime error: " after invoking errorString.Error().
// See Issue #14965.
type plainError string

func (e plainError) RuntimeError() {}

func (e plainError) Error() string {
	return string(e)
}

type stringer interface {
	String() string
}

func typestring(x interface{}) string {
	e := efaceOf(&x)
	return e._type.string()
}

// printany prints an argument passed to panic.
// If panic is called with a value that has a String or Error method,
// it has already been converted into a string by preprintpanics.
func printany(i interface{}) {
	switch v := i.(type) {
	case nil:
		print("nil")
	case bool:
		print(v)
	case int:
		print(v)
	case int8:
		print(v)
	case int16:
		print(v)
	case int32:
		print(v)
	case int64:
		print(v)
	case uint:
		print(v)
	case uint8:
		print(v)
	case uint16:
		print(v)
	case uint32:
		print(v)
	case uint64:
		print(v)
	case uintptr:
		print(v)
	case float32:
		print(v)
	case float64:
		print(v)
	case complex64:
		print(v)
	case complex128:
		print(v)
	case string:
		print(v)
	default:
		print("(", typestring(i), ") ", i)
	}
}

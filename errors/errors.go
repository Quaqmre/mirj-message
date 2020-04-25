package errors

import (
	"bytes"
	"fmt"
)

type Error struct {
	Owner Owner
	Op    Op
	Kind  Kind
	Err   error
}

func (e *Error) isZero() bool {
	return e.Owner == "" && e.Op == "" && e.Kind == 0 && e.Err == nil
}

// Owner is user name
type Owner string

// Op is method and package name holder for error
type Op string

//Kind is grouping separet errors
type Kind uint8

const (
	Other Kind = iota
	Invalid
	NotAccess
	Parse
)

var _ error = (*Error)(nil)

var Separator = ":\n\t"

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case NotAccess:
		return "cant accesible operation"
	case Parse:
		return "parsing error"
	}
	return "unknown error kind"
}

func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case string:
			// Someone might accidentally call us with a user or path name
			// that is not of the right type. Take care of that and log it.
			e.Err = Str(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			// Make a copy
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		default:
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, file name or user name
	// twice.
	if prev.Owner == e.Owner {
		prev.Owner = ""
	}
	if prev.Kind == e.Kind {
		prev.Kind = Other
	}
	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}
	return e
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}
	if e.Owner != "" {
		b.WriteString("owner ")
		b.WriteString(string(e.Owner))
	}
	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Err != nil {
		// Indent on new line if we are cascading non-empty Upspin errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

func Str(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
func Errorf(format string, args ...interface{}) error {
	return &errorString{fmt.Sprintf(format, args...)}
}

func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

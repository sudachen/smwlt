package fu

import (
	"bytes"
	"fmt"
	"golang.org/x/xerrors"
	"strings"
)

type xpanic struct{ err error }

/*
Panic returns traceable panic object
*/
func Panic(err error, skip ...int) interface{} {
	if _, ok := err.(xerrors.Formatter); !ok {
		err = xerror{err, xerrors.Caller(Opti(1, skip...))}
	}
	return xpanic{err}
}

func PanicMessage(e interface{}) string {
	if p, ok := e.(xpanic); ok {
		if err := p.Unwrap(); err != nil {
			return err.Error()
		}
		return p.Error()
	}
	return fmt.Sprint(e)
}

func (x xpanic) stringify(indepth bool) string {
	s, e := stringifyError(x.err)
	ns := []string{s}
	for e != nil && indepth {
		s, e = stringifyError(e)
		ns = append(ns, s)
	}
	return strings.Join(ns, "\n")
}

/*
Error implements error interface
*/
func (x xpanic) Error() string {
	return x.stringify(false)
}

/*
String implements Stringer interface
*/
func (x xpanic) String() string {
	return x.stringify(true)
}

/*
Unwrap implements Wrapper interface
*/
func (x xpanic) Unwrap() error {
	if w, ok := x.err.(xerrors.Wrapper); ok {
		return w.Unwrap()
	}
	return x.err
}

/*
Trace makes error traceable
*/
func Trace(err error) error {
	if _, ok := err.(xerrors.Formatter); ok {
		return err
	}
	return xerror{err, xerrors.Caller(1)}
}

/*
Errorf formats new error
*/
func Errorf(f string, a ...interface{}) error {
	return xerror{fmt.Errorf(f, a...), xerrors.Caller(1)}
}

/*
Wrapf wraps error with formatted string
*/
func Wrapf(err error, f string, a ...interface{}) error {
	return xerror{xwrapper{err, fmt.Sprintf(f, a...)}, xerrors.Caller(1)}
}

/*
Wrap wraps error with string message
*/
func Wrap(err error, m string) error {
	return xerror{xwrapper{err, m}, xerrors.Caller(1)}
}

/*
Error creates new error
*/
func Error(message string) error {
	return xerror{xerrors.New(message), xerrors.Caller(1)}
}

type xerror struct {
	error
	frame xerrors.Frame
}

/*
FormatError implements Formatter interface
*/
func (e xerror) FormatError(p xerrors.Printer) error {
	p.Print(e.error.Error() + " at ")
	e.frame.Format(p)
	return nil
}

func stringifyError(err error) (string, error) {
	ep := &xprinter{details: true}
	if f, ok := err.(xerrors.Formatter); ok {
		err = f.FormatError(ep)
	} else {
		ep.Print(err.Error())
		err = nil
	}
	return strings.Join(strings.Fields(ep.String()), " "), err
}

type xwrapper struct {
	error
	message string
}

/*
Error implements error interface
*/
func (e xwrapper) Error() string {
	return e.message
}

/*
Unwrap implements Wrapper interface
*/
func (e xwrapper) Unwrap() error {
	if w, ok := e.error.(xerrors.Wrapper); ok {
		return w.Unwrap()
	}
	return e.error
}

type xprinter struct {
	bytes.Buffer
	details bool
}

/*
Print implements Printer interface
*/
func (ep *xprinter) Print(args ...interface{}) {
	ep.Buffer.WriteString(fmt.Sprint(args...))
}

/*
Printf implements Printer interface
*/
func (ep *xprinter) Printf(format string, args ...interface{}) {
	ep.Buffer.WriteString(fmt.Sprintf(format, args...))
}

/*
Detail implements Printer interface
*/
func (ep xprinter) Detail() bool {
	return ep.details
}

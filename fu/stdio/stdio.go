package stdio

import (
	"fmt"
	"os"
)

var StdIo = DefaultStdIo

func DefaultStdIo() (stdin, stdout, stderr *os.File) {
	return os.Stdin, os.Stdout, os.Stderr
}

func Print(a ...interface{}) {
	_, stdout, _ := StdIo()
	_, _ = fmt.Fprint(stdout, a...)
}

func Println(a ...interface{}) {
	_, stdout, _ := StdIo()
	_, _ = fmt.Fprintln(stdout, a...)
}

func Printf(f string, a ...interface{}) {
	_, stdout, _ := StdIo()
	_, _ = fmt.Fprintf(stdout, f, a...)
}

func Printfln(f string, a ...interface{}) {
	_, stdout, _ := StdIo()
	_, _ = fmt.Fprintf(stdout, f+"\n", a...)
}

func Errorln(a ...interface{}) {
	_, _, stderr := StdIo()
	_, _ = fmt.Fprintln(stderr, a...)
}

func Errorf(f string, a ...interface{}) {
	_, _, stderr := StdIo()
	_, _ = fmt.Fprintf(stderr, f, a...)
}

func Errorfln(f string, a ...interface{}) {
	_, _, stderr := StdIo()
	_, _ = fmt.Fprintf(stderr, f+"\n", a...)
}

type IoReseter interface{ Reset() }
type reseter []*os.File

func (r reseter) Reset() {
	StdIo = func() (a, b, c *os.File) {
		return r[0], r[1], r[2]
	}
}

func LocalIo(stdin, stdout, stderr *os.File) IoReseter {
	a, b, c := StdIo()
	r := reseter{a, b, c}
	if stdin == nil {
		stdin = a
	}
	if stdout == nil {
		stdout = b
	}
	if stderr == nil {
		stderr = c
	}
	StdIo = func() (a, b, c *os.File) {
		return stdin, stdout, stderr
	}
	return r
}

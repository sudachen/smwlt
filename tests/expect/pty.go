package expect

import (
	"bufio"
	"github.com/sudachen/smwlt/fu/errstr"
	"github.com/sudachen/smwlt/fu/stdio"
	"io"
	"os"
)

type OsIo struct {
	*bufio.Scanner
	Input, Output *os.File
}

func (p OsIo) Close() {
	if p.Input != nil {
		p.Input.Close()
	}
	if p.Output != nil {
		p.Output.Close()
	}
}

func (p OsIo) LocalIo() stdio.IoReseter {
	return stdio.LocalIo(p.Input, p.Output, p.Output)
}

type Pty struct {
	Host   OsIo
	Target OsIo
}

func New() (pty *Pty, err error) {
	pty = &Pty{}
	defer func() {
		if err != nil {
			pty.Close()
		}
	}()
	if pty.Host.Input, pty.Target.Output, err = os.Pipe(); err != nil {
		return
	}
	if pty.Target.Input, pty.Host.Output, err = os.Pipe(); err != nil {
		return
	}

	pty.Host.Scanner = bufio.NewScanner(pty.Host.Input)
	pty.Target.Scanner = bufio.NewScanner(pty.Target.Input)
	return
}

func (pty *Pty) Close() {
	pty.Host.Close()
	pty.Target.Close()
}

func (p OsIo) Send(text string) error {
	_, err := p.Output.WriteString(text)
	return err
}

func (p OsIo) Receive() (string, error) {
	if p.Scan() {
		return p.Text(), nil
	}
	return "", errstr.Wrap(1, io.EOF, "Terminal is Closed")
}

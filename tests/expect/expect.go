package expect

import (
	"github.com/sudachen/smwlt/fu/errstr"
	"regexp"
)

func (p OsIo) ExpectRxOrPanic(rx string) {
	if p.Scan() {
		text := p.Text()
		ok, err := regexp.MatchString(rx, text)
		if err != nil {
			panic(err)
		}
		if !ok {
			panic(errstr.Format(1, "expected `%v`, received `%v`", rx, text))
		}
	}
	return
}

func (p OsIo) SkipRest() {
	for p.Scan() {
	}
}

func (p OsIo) SkipUntil(rx string) {
	for p.Scan() {
		text := p.Text()
		ok, err := regexp.MatchString(rx, text)
		if err != nil {
			panic(err)
		}
		if ok {
			break
		}
	}
}

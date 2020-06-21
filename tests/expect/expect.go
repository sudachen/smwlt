package expect

import (
	"fmt"
	"github.com/sudachen/smwlt/fu/errstr"
	"regexp"
)

func (p OsIo) Expect(rx string){
	reg, err := regexp.Compile(rx)
	if err != nil {
		panic(errstr.Wrapf(1, err, "bad regexp r`%v`: %v", rx, err.Error()))
	}
	if p.Scan() {
		text := p.Text()
		if !reg.MatchString(text) {
			panic(errstr.Format(1, "expected r`%v`, received `%v`", rx, text))
		}
		return
	}
	panic(errstr.Format(1,"not fount r`%v`",rx))
}

func (p OsIo) ExpectGet(rx string) []string{
	reg, err := regexp.Compile(rx)
	if err != nil {
		panic(errstr.Wrapf(1, err, "bad regexp r`%v`: %v", rx, err.Error()))
	}
	if p.Scan() {
		text := p.Text()
		if !reg.MatchString(text) {
			panic(errstr.Format(1, "expected r`%v`, received `%v`", rx, text))
		}
		return reg.FindStringSubmatch(text)[1:]
	}
	panic(errstr.Format(1,"not fount r`%v`",rx))
}

func (p OsIo) SkipRest() {
	for p.Scan() {
	}
}

func (p OsIo) SkipToExpect(rx string) {
	reg, err := regexp.Compile(rx)
	if err != nil {
		panic(errstr.Wrapf(1, err, "bad regexp r`%v`: %v", rx, err.Error()))
	}
	for p.Scan() {
		text := p.Text()
		if reg.MatchString(text) {
			return
		}
		fmt.Println(text)
	}
	panic(errstr.Format(1,"not fount r`%v`",rx))
}

func (p OsIo) SkipToExpectGet(rx string) []string{
	reg, err := regexp.Compile(rx)
	if err != nil {
		panic(errstr.Wrapf(1, err, "bad regexp r`%v`: %v", rx, err.Error()))
	}
	for p.Scan() {
		text := p.Text()
		if reg.MatchString(text) {
			return reg.FindStringSubmatch(text)[1:]
		}
	}
	panic(errstr.Format(1,"not fount r`%v`",rx))
}

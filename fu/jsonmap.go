package fu

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sudachen/smwlt/fu/errstr"
	"io"
)

var JsonTypeError = errors.New("JsonMap value type error")

type JsonPath string

type JsonMap struct {
	Val  map[string]interface{}
	path JsonPath
}

type JsonList struct {
	Val  []interface{}
	path JsonPath
}

type JsonValue struct {
	Val  interface{}
	path JsonPath
}

func (p JsonPath) Next(n string) JsonPath {
	if p == "" {
		return JsonPath(n)
	}
	return p + JsonPath("."+n)
}

func (m *JsonMap) Decode(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&m.Val)
}

func (m *JsonMap) Encode(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(&m.Val)
}

func (m JsonMap) Map(n string) JsonMap {
	next := m.path.Next(n)
	if v, ok := m.Val[n]; ok {
		if q, ok := v.(map[string]interface{}); ok {
			return JsonMap{q, next}
		}
		panic(errstr.Wrapf(1, JsonTypeError, "value '%v' is not a map", next))
	}
	return JsonMap{map[string]interface{}{}, next}
}

func (m JsonMap) Len() int {
	return len(m.Val)
}

func (m JsonMap) Value(n string) JsonValue {
	next := m.path.Next(n)
	if v, ok := m.Val[n]; ok {
		return JsonValue{v, next}
	}
	return JsonValue{"", next}
}

func (m JsonMap) List(n string) JsonList {
	next := m.path.Next(n)
	if v, ok := m.Val[n]; ok {
		if q, ok := v.([]interface{}); ok {
			return JsonList{q, next}
		}
		panic(errstr.Wrapf(1,JsonTypeError, "value '%v' is not a list", next))
	}
	return JsonList{[]interface{}{}, next}
}

func (l JsonList) Len() int {
	return len(l.Val)
}

func (l JsonList) Value(i int) JsonValue {
	return JsonValue{l.Val[i], l.path.Next(fmt.Sprintf("[%v]", i))}
}

func (l JsonList) Maps() []JsonMap {
	r := make([]JsonMap, l.Len())
	for i, v := range l.Val {
		next := l.path.Next(fmt.Sprintf("[%v]", i))
		if x, ok := v.(map[string]interface{}); ok {
			r[i] = JsonMap{x, next}
		} else {
			panic(errstr.Wrapf(1,JsonTypeError, "value '%v' is not a map", next))
		}
	}
	return r
}

func (v JsonValue) String() string {
	if q, ok := v.Val.(string); ok {
		return q
	}
	panic(errstr.Wrapf(1, JsonTypeError, "value '%v' is not a string", v.path))
}

func (v JsonValue) HexBytes() []byte {
	s := v.String()
	if s != "" {
		bs, err := hex.DecodeString(s)
		if err != nil {
			panic(errstr.Wrapf(1,JsonTypeError, "value '%v' is not a hex string", v.path))
		}
		return bs
	}
	return nil
}

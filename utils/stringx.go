package utils

import (
	"encoding/json"
	"fmt"
)

type Stringer struct {
	v interface{}
}

func NewStringer(v interface{}) *Stringer {
	var stringer = new(Stringer)
	stringer.v = v
	return stringer
}

func (stringer Stringer) String() string {
	if stringer.v == nil {
		return ""
	}
	switch stringer.v.(type) {
	case string:
		return stringer.v.(string)
	case []rune:
		return string(stringer.v.([]rune))
	case []byte:
		return string(stringer.v.([]byte))
	case fmt.Stringer:
		return stringer.v.(fmt.Stringer).String()
	case fmt.GoStringer:
		return stringer.v.(fmt.GoStringer).GoString()
	case json.Marshaler:
		var bytes, err = stringer.v.(json.Marshaler).MarshalJSON()
		if err == nil {
			return string(bytes)
		}
	}
	return fmt.Sprintf("%v", stringer.v)
}

func (stringer Stringer) GoString() string {
	return stringer.String()
}

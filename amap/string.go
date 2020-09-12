package amap

import (
	"bytes"
	"errors"
	"github.com/brickman-source/golang-utilities/json"
)

var (
	sbArray = []byte(`[]`)
)
type String []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m String) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte(`""`), nil
	}
	if bytes.Compare(sbArray, m) == 0 {
		return []byte(`""`), nil
	}
	if len(m) == 0 {
		return []byte(`""`), nil
	}
	return json.Marshal(string(m))
}

// UnmarshalJSON sets *m to a copy of data.
func (m *String) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("amap.String: UnmarshalJSON on nil pointer")
	}
	if bytes.Compare(sbArray, data) == 0 {
		*m = []byte(``)
		return nil
	}
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*m = []byte(str) //
	return nil
}

func (m String) ToString() string {
	if m == nil {
		println("m == nil")
		return ""
	}
	if bytes.Compare(sbArray, m) == 0 {
		println("m is sb array")
		return ""
	}
	if len(m) < 1 {
		println("len(m) < 1")
		return ""
	}
	return string(m)
}

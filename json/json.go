package json

import (
	"bytes"
	"github.com/json-iterator/go"
)

func MarshalToString(v interface{}) (string, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(v)
}

func MarshalToBuffer(v interface{}) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(buf).Encode(v)
	return buf, err
}

func UnmarshalFromString(str string, v interface{}) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(str, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}

func ShouldMarshal(v interface{}) []byte {
	ret, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
	if err != nil {
		return nil
	}
	return ret
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, v)
}

package number

import (
	"encoding/json"
	"fmt"
	//"github.com/mgo/bson"
)

type Float64p2 float64

func (n Float64p2) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%0.2f", n)), nil
}

func (n *Float64p2) UnmarshalJSON(b []byte) error {
	var f float64
	err := json.Unmarshal(b, &f)
	*n = Float64p2(f)
	return err
}

func (n Float64p2) String() string {
	return fmt.Sprintf("%0.2f", n)
}

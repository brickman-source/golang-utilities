package datetime

import (
	"github.com/araddon/dateparse"
	"github.com/brickman-source/golang-utilities/json"
	"strings"
	"time"
)

type UTCTime time.Time

func (dt UTCTime) MarshalJSON() ([]byte, error) {
	tmp := time.Time(dt)
	return json.ShouldMarshal(tmp.UTC().Format(time.RFC3339)), nil
}

func (dt *UTCTime) UnmarshalJSON(p []byte) error {
	timeStr := strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	)
	t, err := dateparse.ParseAny(timeStr)
	if err != nil {
		t, err = time.Parse(time.RFC3339, strings.Replace(
			string(p),
			"\"",
			"",
			-1,
		))
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04Z07:00", strings.Replace(
				string(p),
				"\"",
				"",
				-1,
			))
		}
	}
	if err != nil {
		return err
	}
	*dt = UTCTime(t)

	return nil
}

func (dt *UTCTime) BeforeEpoch() bool {
	if dt == nil {
		return false
	}
	t := time.Time(*dt)
	return t.Before(time.Unix(0, 0))
}

func (dt *UTCTime) AfterEpoch() bool {
	if dt == nil {
		return false
	}
	t := time.Time(*dt)
	return t.After(time.Unix(0, 0))
}

func (dt *UTCTime) Before(u time.Time) bool {
	if dt == nil {
		return false
	}
	t := time.Time(*dt)
	return t.Before(u)
}

func (dt *UTCTime) After(u time.Time) bool {
	if dt == nil {
		return false
	}
	t := time.Time(*dt)
	return t.After(u)
}

func (dt *UTCTime) Unix() int64 {
	if dt == nil {
		return 0
	}
	t := time.Time(*dt)
	return t.Unix()
}

func Unix(sec int64, nsec int64) UTCTime {
	t := time.Unix(sec, nsec)
	return UTCTime(t)
}

func (dt *UTCTime) Format(layout string) string  {
	if dt == nil {
		return ""
	}
	t := time.Time(*dt)
	return t.Format(layout)
}

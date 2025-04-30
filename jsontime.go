package certmon

import (
	"fmt"
	"strings"
	"time"
)

type JSONTime time.Time

const timeLayout = time.DateTime

func (jt JSONTime) MarshalJSON() ([]byte, error) {
	t := time.Time(jt)
	s := fmt.Sprintf(`"%s"`, t.Format(timeLayout))
	return []byte(s), nil
}

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(timeLayout, s)
	if err != nil {
		return err
	}
	*jt = JSONTime(t)
	return nil
}

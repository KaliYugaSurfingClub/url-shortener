package valkit

import (
	"fmt"
	"github.com/thoas/go-funk"
	"strconv"
	"strings"
)

func addBrackets(s string) string {
	return "'" + s + "'"
}

func ContainsInMap[V any](acceptable map[string]V) func(value any) error {
	keys := funk.Keys(acceptable).([]string)
	inBuckets := strings.Join(funk.Map(keys, addBrackets).([]string), ", ")

	return func(value any) error {
		if _, ok := acceptable[value.(string)]; !ok {
			return fmt.Errorf("should be one from: %s", strings.TrimSpace(inBuckets))
		}

		return nil
	}
}

func StringNumIn(min int64, max int64) func(value any) error {
	return func(value any) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("internal error: invalid string")
		}

		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil || num < min || num > max {
			return fmt.Errorf("should be a integer in [%d,%d]", min, max)
		}

		return nil
	}
}

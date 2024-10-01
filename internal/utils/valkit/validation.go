package valkit

import (
	"fmt"
	"github.com/thoas/go-funk"
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

func Positive() func(value any) error {
	return func(value any) error {
		if value.(int64) <= 0 {
			return fmt.Errorf("should be positive")
		}

		return nil
	}
}

package agent

import (
	"fmt"
)

func GetFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func GetInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int64:
		return v, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

package utility

import (
	"fmt"
	"strconv"
)

type StrMap = map[string]interface{}
type AnyMap = map[interface{}]interface{}

func AnyToAnyMap(value interface{}) AnyMap {
	if value == nil {
		return nil
	}
	switch val := value.(type) {
	case AnyMap:
		return val
	case StrMap:
		count := len(val)
		if count == 0 {
			return nil
		}
		m := make(AnyMap, count)
		for k, v := range val {
			m[k] = v
		}
		return m
	default:
		return nil
	}
}

func AnyToStrMap(value interface{}) StrMap {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case StrMap:
		return v
	case AnyMap:
		l := len(v)
		if l == 0 {
			return nil
		}
		m := make(StrMap, l)
		for k, v := range v {
			m[AnyToString(k)] = v
		}
		return m
	default:
		return nil
	}
}

func AnyToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch val := value.(type) {
	case *string:
		if val == nil {
			return ""
		}
		return *val
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case error:
		return val.Error()
	default:
		return fmt.Sprint(value)
	}
}

func AnyToStringArray(any interface{}) []string {
	if any == nil {
		return nil
	}
	switch v := any.(type) {
	case []string:
		return v
	case []interface{}:
		return AnyArrayToStringArray(v)
	default:
		return nil
	}
}

func AnyArrayToStringArray(arrInterface []interface{}) []string {
	elementArray := make([]string, len(arrInterface))
	for i, v := range arrInterface {
		elementArray[i] = AnyToString(v)
	}
	return elementArray
}

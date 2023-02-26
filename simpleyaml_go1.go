package simpleyaml

import (
	"errors"
	"reflect"
)

// Float64 coerces into a float64
func (y *Yaml) Float64() (float64, error) {
	switch y.data.(type) {
	case float32, float64:
		return reflect.ValueOf(y.data).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(y.data).Int()), nil
	}
	return 0, errors.New("invalid value type")
}

// Int coerces into an int
func (y *Yaml) Int() (int, error) {
	switch y.data.(type) {
	case float32, float64:
		return int(reflect.ValueOf(y.data).Float()), nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(y.data).Int()), nil
	}
	return 0, errors.New("invalid value type")
}

// Int64 coerces into an int64
func (y *Yaml) Int64() (int64, error) {
	switch y.data.(type) {
	case float32, float64:
		return int64(reflect.ValueOf(y.data).Float()), nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(y.data).Int(), nil
	}
	return 0, errors.New("invalid value type")
}

// Uint64 coerces into an uint64
func (y *Yaml) Uint64() (uint64, error) {
	switch y.data.(type) {
	case float32, float64:
		return uint64(reflect.ValueOf(y.data).Float()), nil
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(y.data).Int()), nil
	}
	return 0, errors.New("invalid value type")
}

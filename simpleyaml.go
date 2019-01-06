package simpleyaml

import (
	"errors"
	"io"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// Version returns the current implementation version
func Version() string {
	return "0.1.0"
}

// Yaml Yaml data
type Yaml struct {
	data interface{}
}

// New returns a pointer to a new, empty `Yaml` object
func New() *Yaml {
	return &Yaml{
		data: map[interface{}]interface{}{},
	}
}

// NewFromReader returns a *Yaml by decoding from an io.Reader
func NewFromReader(r io.Reader) (*Yaml, error) {
	y := new(Yaml)
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&y.data)
	return y, err
}

// NewYaml returns a pointer to a new `Yaml` object
// after unmarshaling `body` bytes
func NewYaml(body []byte) (*Yaml, error) {
	y := new(Yaml)
	err := y.Unmarshal(body)
	if err != nil {
		return nil, err
	}
	return y, nil
}

// Interface returns the underlying data
func (y *Yaml) Interface() interface{} {
	return y.data
}

// Marshal Implements the yaml.Marshaler interface.
func (y *Yaml) Marshal() ([]byte, error) {
	return yaml.Marshal(&y.data)
}

// Unmarshal Implements the yaml.Unmarshaler interface.
func (y *Yaml) Unmarshal(p []byte) error {
	return yaml.Unmarshal(p, &y.data)
}

// Set modifies `Yaml` map by `key` and `value`
// Useful for changing single key/value in a `Yaml` object easily.
func (y *Yaml) Set(key string, val interface{}) {
	m, err := y.Map()
	if err != nil {
		return
	}
	m[key] = val
}

// SetPath modifies `Yaml`, recursively checking/creating map keys for the supplied path,
// and then finally writing in the value
func (y *Yaml) SetPath(branch []string, val interface{}) {
	if len(branch) == 0 {
		y.data = val
		return
	}

	// in order to insert our branch, we need map[interface{}]interface{}
	if _, ok := (y.data).(map[interface{}]interface{}); !ok {
		// have to replace with something suitable
		y.data = make(map[interface{}]interface{})
	}
	curr := y.data.(map[interface{}]interface{})

	for i := 0; i < len(branch)-1; i++ {
		b := branch[i]
		// key exists?
		if _, ok := curr[b]; !ok {
			n := make(map[interface{}]interface{})
			curr[b] = n
			curr = n
			continue
		}

		// make sure the value is the right sort of thing
		if _, ok := curr[b].(map[interface{}]interface{}); !ok {
			// have to replace with something suitable
			n := make(map[interface{}]interface{})
			curr[b] = n
		}

		curr = curr[b].(map[interface{}]interface{})
	}

	// add remaining k/v
	curr[branch[len(branch)-1]] = val
}

// Del modifies `Yaml` map by deleting `key` if it is present.
func (y *Yaml) Del(key string) {
	m, err := y.Map()
	if err != nil {
		return
	}
	delete(m, key)
}

// Get returns a pointer to a new `Yaml` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested YAML):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (y *Yaml) Get(key string) *Yaml {
	m, err := y.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Yaml{val}
		}
	}
	return &Yaml{nil}
}

// GetPath searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetPath("top_level", "dict")
func (y *Yaml) GetPath(branch ...string) *Yaml {
	jin := y
	for _, p := range branch {
		jin = jin.Get(p)
	}
	return jin
}

// GetIndex returns a pointer to a new `Yaml` object
// for `index` in its `array` representation
//
// this is the analog to Get when accessing elements of
// a array instead of a object:
//    js.Get("top_level").Get("array").GetIndex(1).Get("key").Int()
func (y *Yaml) GetIndex(index int) *Yaml {
	a, err := y.Array()
	if err == nil {
		if len(a) > index {
			return &Yaml{a[index]}
		}
	}
	return &Yaml{nil}
}

// CheckGet returns a pointer to a new `Yaml` object and
// a `bool` identifying success or failure
//
// useful for chained operations when success is important:
//    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
//        log.Println(data)
//    }
func (y *Yaml) CheckGet(key string) (*Yaml, bool) {
	m, err := y.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Yaml{val}, true
		}
	}
	return nil, false
}

// Map type asserts to `map`
func (y *Yaml) Map() (map[interface{}]interface{}, error) {
	if m, ok := (y.data).(map[interface{}]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[interface{}]interface{} failed")
}

// Array type asserts to an `array`
func (y *Yaml) Array() ([]interface{}, error) {
	if a, ok := (y.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

// Bool type asserts to `bool`
func (y *Yaml) Bool() (bool, error) {
	if s, ok := (y.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `string`
func (y *Yaml) String() (string, error) {
	if s, ok := (y.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// Bytes type asserts to `[]byte`
func (y *Yaml) Bytes() ([]byte, error) {
	if s, ok := (y.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}

// StringArray type asserts to an `array` of `string`
func (y *Yaml) StringArray() ([]string, error) {
	arr, err := y.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, "")
			continue
		}
		s, ok := a.(string)
		if !ok {
			return nil, errors.New("type assertion to []string failed")
		}
		retArr = append(retArr, s)
	}
	return retArr, nil
}

// MustArray guarantees the return of a `[]interface{}` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, v := range js.Get("results").MustArray() {
//			fmt.Println(i, v)
//		}
func (y *Yaml) MustArray(args ...[]interface{}) []interface{} {
	var def []interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustArray() received too many arguments %d", len(args))
	}

	a, err := y.Array()
	if err == nil {
		return a
	}

	return def
}

// MustMap guarantees the return of a `map[interface{}]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").MustMap() {
//			fmt.Println(k, v)
//		}
func (y *Yaml) MustMap(args ...map[interface{}]interface{}) map[interface{}]interface{} {
	var def map[interface{}]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustMap() received too many arguments %d", len(args))
	}

	a, err := y.Map()
	if err == nil {
		return a
	}

	return def
}

// MustString guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").MustString(), js.Get("optional_param").MustString("my_default"))
func (y *Yaml) MustString(args ...string) string {
	var def string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustString() received too many arguments %d", len(args))
	}

	s, err := y.String()
	if err == nil {
		return s
	}

	return def
}

// MustStringArray guarantees the return of a `[]string` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, s := range js.Get("results").MustStringArray() {
//			fmt.Println(i, s)
//		}
func (y *Yaml) MustStringArray(args ...[]string) []string {
	var def []string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustStringArray() received too many arguments %d", len(args))
	}

	a, err := y.StringArray()
	if err == nil {
		return a
	}

	return def
}

// MustInt guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.Get("param1").MustInt(), js.Get("optional_param").MustInt(5150))
func (y *Yaml) MustInt(args ...int) int {
	var def int

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt() received too many arguments %d", len(args))
	}

	i, err := y.Int()
	if err == nil {
		return i
	}

	return def
}

// MustFloat64 guarantees the return of a `float64` (with optional default)
//
// useful when you explicitly want a `float64` in a single value return context:
//     myFunc(js.Get("param1").MustFloat64(), js.Get("optional_param").MustFloat64(5.150))
func (y *Yaml) MustFloat64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustFloat64() received too many arguments %d", len(args))
	}

	f, err := y.Float64()
	if err == nil {
		return f
	}

	return def
}

// MustBool guarantees the return of a `bool` (with optional default)
//
// useful when you explicitly want a `bool` in a single value return context:
//     myFunc(js.Get("param1").MustBool(), js.Get("optional_param").MustBool(true))
func (y *Yaml) MustBool(args ...bool) bool {
	var def bool

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustBool() received too many arguments %d", len(args))
	}

	b, err := y.Bool()
	if err == nil {
		return b
	}

	return def
}

// MustInt64 guarantees the return of an `int64` (with optional default)
//
// useful when you explicitly want an `int64` in a single value return context:
//     myFunc(js.Get("param1").MustInt64(), js.Get("optional_param").MustInt64(5150))
func (y *Yaml) MustInt64(args ...int64) int64 {
	var def int64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt64() received too many arguments %d", len(args))
	}

	i, err := y.Int64()
	if err == nil {
		return i
	}

	return def
}

// MustUint64 guarantees the return of an `uint64` (with optional default)
//
// useful when you explicitly want an `uint64` in a single value return context:
//     myFunc(js.Get("param1").MustUint64(), js.Get("optional_param").MustUint64(5150))
func (y *Yaml) MustUint64(args ...uint64) uint64 {
	var def uint64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustUint64() received too many arguments %d", len(args))
	}

	i, err := y.Uint64()
	if err == nil {
		return i
	}

	return def
}

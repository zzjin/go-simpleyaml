package simpleyaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleYAMLGo1(t *testing.T) {
	sy, err := NewYaml([]byte(`
test:
  float64: 30.02
  int64: 8000000000000000008
  int: -9527
  string: "simpleyaml"
`))
	assert.NotEqual(t, sy, nil)
	if !assert.Nil(t, err) {
		t.Fail()
	}

	arr, _ := sy.Get("test").Array()
	assert.NotEqual(t, arr, nil)

	mm := sy.Get("test").MustMap()
	assert.Equal(t, mm, map[interface{}]interface{}{
		"float64": float64(30.02),
		"int64":   int(8000000000000000008), // belows are all default to int
		"int":     int(-9527),
		"string":  "simpleyaml",
	})

	ne := sy.Get("test").Get("not_exists")
	assert.Empty(t, ne)

	//for Float64
	var f64 float64
	f64, _ = sy.Get("test").Get("float64").Float64()
	assert.Equal(t, f64, float64(30.02))
	f64, _ = sy.Get("test").Get("int64").Float64()
	assert.Equal(t, f64, float64(8000000000000000008))
	f64, err = sy.Get("test").Get("string").Float64()
	assert.NotNil(t, err)

	//for Int
	var defaultInt int
	defaultInt, _ = sy.Get("test").Get("float64").Int()
	assert.Equal(t, defaultInt, int(30))
	defaultInt, _ = sy.Get("test").Get("int64").Int()
	assert.Equal(t, defaultInt, int(8000000000000000008))
	defaultInt, _ = sy.Get("test").Get("int").Int()
	assert.Equal(t, defaultInt, int(-9527))
	defaultInt, err = sy.Get("test").Get("string").Int()
	assert.NotNil(t, err)

	//for Int64
	var defaultInt64 int64
	defaultInt64, _ = sy.Get("test").Get("float64").Int64()
	assert.Equal(t, defaultInt64, int64(30))
	defaultInt64, _ = sy.Get("test").Get("int64").Int64()
	assert.Equal(t, defaultInt64, int64(8000000000000000008))
	defaultInt64, err = sy.Get("test").Get("string").Int64()
	assert.NotNil(t, err)

	//for Uint64
	var defaultUint64 uint64
	defaultUint64, _ = sy.Get("test").Get("float64").Uint64()
	assert.Equal(t, defaultUint64, uint64(30))
	defaultUint64, _ = sy.Get("test").Get("int64").Uint64()
	assert.Equal(t, defaultUint64, uint64(8000000000000000008))
	defaultUint64, _ = sy.Get("test").Get("int").Uint64()
	assert.Equal(t, defaultUint64, uint64(18446744073709542089))
	defaultUint64, err = sy.Get("test").Get("string").Uint64()
	assert.NotNil(t, err)
}

package simpleyaml

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleYAMLGo1(t *testing.T) {
	sy, newErr := NewYaml([]byte(`
test:
  array: [1, "2", 3.0]
  arraywithsubs:
    - subkeyone: 1
    - subkeytwo: 2
      subkeythree: 3
  bignum: 800000000000000008
`))
	assert.NotEqual(t, nil, sy)
	assert.Equal(t, nil, newErr)

	arr, _ := sy.Get("test").Get("array").Array()
	assert.NotEqual(t, nil, arr)
	for i, v := range arr {
		var iv int
		switch v.(type) {
		case int:
			iv = v.(int)
		case float64:
			iv = int(v.(float64))
		case string:
			iv, _ = strconv.Atoi(v.(string))
		}
		assert.Equal(t, i+1, iv)
	}

	ma := sy.Get("test").Get("array").MustArray()
	assert.Equal(t, ma, []interface{}{int(1), "2", float64(3)})

	mm := sy.Get("test").Get("arraywithsubs").GetIndex(0).MustMap()
	assert.Equal(t, mm, map[interface{}]interface{}{"subkeyone": int(1)})
	assert.Equal(t, mm["subkeyone"], int(1))

	bigint := sy.Get("test").Get("bignum").MustInt64()
	assert.Equal(t, bigint, int64(800000000000000008))
}

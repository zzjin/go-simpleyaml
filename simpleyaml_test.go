package simpleyaml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.NotEmpty(t, Version())
}

func TestSimpleYAML(t *testing.T) {
	var ok bool
	var err error
	var sy *Yaml

	sy = New()
	assert.Empty(t, sy.Interface())

	sy.SetPath([]string{}, "test")
	assert.NotEmpty(t, sy.Interface())

	sy.SetPath([]string{"test", "string"}, "simpleyaml")
	assert.NotEmpty(t, sy.Interface())

	sy.SetPath([]string{"test", "string"}, "simpleyaml2")
	assert.NotEmpty(t, sy.Interface())

	//tab to make err
	sy, err = NewYaml([]byte(`
	tab
`))
	assert.NotNil(t, err)

	testData := []byte(`
test:
  string_array: ["asdf", "ghjk", "zxcv"]
  string_array_null: ["abc", null, "efg"]
  array: [1, "2", 3.0]
  arraywithsubs:
    - subkeyone: 1
    - subkeytwo: 2
      subkeythree: 3
  int: 10
  float: 5.150
  string: simplesyaml
  bool: true
  sub_obj:
    a: 1
`)

	reader := bytes.NewReader(testData)
	sy, err = NewFromReader(reader)
	assert.NotEqual(t, nil, sy)
	assert.Nil(t, err)

	sy, err = NewYaml(testData)
	assert.NotEqual(t, nil, sy)
	assert.Nil(t, err)

	assert.NotEmpty(t, sy.Interface())

	encodeBytes, err := sy.Marshal()
	assert.Len(t, encodeBytes, 273)

	_, ok = sy.CheckGet("test")
	assert.Equal(t, true, ok)

	_, ok = sy.CheckGet("missing_key")
	assert.Equal(t, false, ok)

	aws := sy.Get("test").Get("arraywithsubs")
	assert.NotEqual(t, nil, aws)
	var awsval int
	awsval, _ = aws.GetIndex(0).Get("subkeyone").Int()
	assert.Equal(t, 1, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeytwo").Int()
	assert.Equal(t, 2, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeythree").Int()
	assert.Equal(t, 3, awsval)
	awsnot := aws.GetIndex(2)
	assert.Equal(t, nil, awsnot.Interface())

	i, _ := sy.Get("test").Get("int").Int()
	assert.Equal(t, 10, i)

	f, _ := sy.Get("test").Get("float").Float64()
	assert.Equal(t, 5.150, f)

	s, _ := sy.Get("test").Get("string").String()
	assert.Equal(t, "simplesyaml", s)

	sb, _ := sy.Get("test").Get("string").Bytes()
	assert.Equal(t, []byte("simplesyaml"), sb)

	_, err = sy.Get("test").Get("string_array").Bytes()
	assert.NotNil(t, err)

	b, _ := sy.Get("test").Get("bool").Bool()
	assert.Equal(t, true, b)
	_, err = sy.Get("test").Get("int").Bool()
	assert.NotNil(t, err)

	mi := sy.Get("test").Get("int").MustInt()
	assert.Equal(t, 10, mi)

	mi64 := sy.Get("test").Get("int").MustInt64()
	assert.Equal(t, int64(10), mi64)
	mi64Not := sy.Get("test").Get("int_not").MustInt64(10)
	assert.Equal(t, int64(10), mi64Not)

	mui := sy.Get("test").Get("int").MustUint64()
	assert.Equal(t, uint64(10), mui)
	muiNot := sy.Get("test").Get("int_not").MustUint64(10)
	assert.Equal(t, uint64(10), muiNot)

	mi2 := sy.Get("test").Get("missing_int").MustInt(5150)
	assert.Equal(t, 5150, mi2)

	ms := sy.Get("test").Get("string").MustString()
	assert.Equal(t, "simplesyaml", ms)

	ms2 := sy.Get("test").Get("missing_string").MustString("fyea")
	assert.Equal(t, "fyea", ms2)

	ma2 := sy.Get("test").Get("missing_array").MustArray([]interface{}{"1", 2, "3"})
	assert.Equal(t, ma2, []interface{}{"1", 2, "3"})

	msa := sy.Get("test").Get("string_array").MustStringArray()
	assert.Equal(t, msa[0], "asdf")
	assert.Equal(t, msa[1], "ghjk")
	assert.Equal(t, msa[2], "zxcv")

	msa2 := sy.Get("test").Get("string_array").MustStringArray([]string{"1", "2", "3"})
	assert.Equal(t, msa2[0], "asdf")
	assert.Equal(t, msa2[1], "ghjk")
	assert.Equal(t, msa2[2], "zxcv")

	msa3 := sy.Get("test").Get("missing_array").MustStringArray([]string{"1", "2", "3"})
	assert.Equal(t, msa3, []string{"1", "2", "3"})

	mm2 := sy.Get("test").Get("missing_map").MustMap(map[interface{}]interface{}{"found": false})
	assert.Equal(t, mm2, map[interface{}]interface{}{"found": false})

	strs, err := sy.Get("test").Get("string_array").StringArray()
	assert.Equal(t, err, nil)
	assert.Equal(t, strs[0], "asdf")
	assert.Equal(t, strs[1], "ghjk")
	assert.Equal(t, strs[2], "zxcv")

	strs2, err := sy.Get("test").Get("string_array_null").StringArray()
	assert.Equal(t, err, nil)
	assert.Equal(t, strs2[0], "abc")
	assert.Equal(t, strs2[1], "")
	assert.Equal(t, strs2[2], "efg")

	gp, _ := sy.GetPath("test", "string").String()
	assert.Equal(t, "simplesyaml", gp)

	gp2, _ := sy.GetPath("test", "int").Int()
	assert.Equal(t, 10, gp2)

	assert.Equal(t, sy.Get("test").Get("bool").MustBool(), true)
	assert.Equal(t, sy.Get("test").Get("bool_not").MustBool(false), false)

	sy.Set("float2", 300.0)
	assert.Equal(t, sy.Get("float2").MustFloat64(), 300.0)
	assert.Equal(t, sy.Get("float_no").MustFloat64(200.0), 200.0)

	sy.Set("test2", "setTest")
	assert.Equal(t, "setTest", sy.Get("test2").MustString())

	sy.Del("test2")
	assert.NotEqual(t, "setTest", sy.Get("test2").MustString())

	sy.Get("test").Get("sub_obj").Set("a", 2)
	assert.Equal(t, 2, sy.Get("test").Get("sub_obj").Get("a").MustInt())

	sy.GetPath("test", "sub_obj").Set("a", 3)
	assert.Equal(t, 3, sy.GetPath("test", "sub_obj", "a").MustInt())
}

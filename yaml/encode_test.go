// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yaml

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

var tMarshalTests = []struct {
	value interface{}
	data  string
}{
	{
		&struct{}{},
		"{}\n",
	}, {
		map[string]string{"v": "hi"},
		"v: hi\n",
	}, {
		map[string]interface{}{"v": "hi"},
		"v: hi\n",
	}, {
		map[string]string{"v": "true"},
		"v: \"true\"\n",
	}, {
		map[string]string{"v": "false"},
		"v: \"false\"\n",
	}, {
		map[string]interface{}{"v": true},
		"v: true\n",
	}, {
		map[string]interface{}{"v": false},
		"v: false\n",
	}, {
		map[string]interface{}{"v": 10},
		"v: 10\n",
	}, {
		map[string]interface{}{"v": -10},
		"v: -10\n",
	}, {
		map[string]uint{"v": 42},
		"v: 42\n",
	}, {
		map[string]interface{}{"v": int64(4294967296)},
		"v: 4294967296\n",
	}, {
		map[string]int64{"v": int64(4294967296)},
		"v: 4294967296\n",
	}, {
		map[string]uint64{"v": 4294967296},
		"v: 4294967296\n",
	}, {
		map[string]interface{}{"v": "10"},
		"v: \"10\"\n",
	}, {
		map[string]interface{}{"v": 0.1},
		"v: 0.1\n",
	}, {
		map[string]interface{}{"v": float64(0.1)},
		"v: 0.1\n",
	}, {
		map[string]interface{}{"v": -0.1},
		"v: -0.1\n",
	}, {
		map[string]interface{}{"v": math.Inf(+1)},
		"v: .inf\n",
	}, {
		map[string]interface{}{"v": math.Inf(-1)},
		"v: -.inf\n",
	}, {
		map[string]interface{}{"v": math.NaN()},
		"v: .nan\n",
	}, {
		map[string]interface{}{"v": nil},
		"v: null\n",
	}, {
		map[string]interface{}{"v": ""},
		"v: \"\"\n",
	}, {
		map[string][]string{"v": []string{"A", "B"}},
		"v:\n- A\n- B\n",
	}, {
		map[string][]string{"v": []string{"A", "B\nC"}},
		"v:\n- A\n- 'B\n\n  C'\n",
	}, {
		map[string][]interface{}{"v": []interface{}{"A", 1, map[string][]int{"B": []int{2, 3}}}},
		"v:\n- A\n- 1\n- B:\n  - 2\n  - 3\n",
	}, {
		map[string]interface{}{"a": map[interface{}]interface{}{"b": "c"}},
		"a:\n  b: c\n",
	}, {
		map[string]interface{}{"a": "-"},
		"a: '-'\n",
	},

	// Simple values.
	{
		tNewInt(123),
		"123\n",
	},

	// Structures
	{
		&struct{ Hello string }{"world"},
		"hello: world\n",
	}, {
		&struct {
			A struct {
				B string
			}
		}{struct{ B string }{"c"}},
		"a:\n  b: c\n",
	}, {
		&struct {
			A *struct {
				B string
			}
		}{&struct{ B string }{"c"}},
		"a:\n  b: c\n",
	}, {
		&struct {
			A *struct {
				B string
			}
		}{},
		"a: null\n",
	}, {
		&struct{ A int }{1},
		"a: 1\n",
	}, {
		&struct{ A []int }{[]int{1, 2}},
		"a:\n- 1\n- 2\n",
	}, {
		&struct {
			B int "a"
		}{1},
		"a: 1\n",
	}, {
		&struct{ A bool }{true},
		"a: true\n",
	},

	// Conditional flag
	{
		&struct {
			A int "a,omitempty"
			B int "b,omitempty"
		}{1, 0},
		"a: 1\n",
	}, {
		&struct {
			A int "a,omitempty"
			B int "b,omitempty"
		}{0, 0},
		"{}\n",
	}, {
		&struct {
			A *struct{ X int } "a,omitempty"
			B int              "b,omitempty"
		}{nil, 0},
		"{}\n",
	},

	// Flow flag
	{
		&struct {
			A []int "a,flow"
		}{[]int{1, 2}},
		"a: [1, 2]\n",
	}, {
		&struct {
			A map[string]string "a,flow"
		}{map[string]string{"b": "c", "d": "e"}},
		"a: {b: c, d: e}\n",
	}, {
		&struct {
			A struct {
				B, D string
			} "a,flow"
		}{struct{ B, D string }{"c", "e"}},
		"a: {b: c, d: e}\n",
	},

	// Unexported field
	{
		&struct {
			u int
			A int
		}{0, 1},
		"a: 1\n",
	},

	// Ignored field
	{
		&struct {
			A int
			B int "-"
		}{1, 2},
		"a: 1\n",
	},

	// Struct inlining
	{
		&struct {
			A int
			C tInlineB `yaml:",inline"`
		}{1, tInlineB{2, tInlineC{3}}},
		"a: 1\nb: 2\nc: 3\n",
	},

	// Duration
	{
		map[string]time.Duration{"a": 3 * time.Second},
		"a: 3s\n",
	},

	// Issue #24.
	{
		map[string]string{"a": "<foo>"},
		"a: <foo>\n",
	},
}

func TestMarshal(t *testing.T) {
	for i, item := range tMarshalTests {
		data, err := Marshal(item.value)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		if s := string(data); s != item.data {
			t.Fatalf("%d: expect = %v, got = %c", i, item.data, s)
		}
	}
}

var tMarshalErrorTests = []struct {
	value interface{}
	error string
}{
	{
		&struct {
			B        int
			tInlineB ",inline"
		}{1, tInlineB{2, tInlineC{3}}},
		`Duplicated key 'b' in struct struct \{ B int; .*`,
	},
}

func TestMarshalErrors(t *testing.T) {
	for i, item := range tMarshalErrorTests {
		_, err := Marshal(item.value)
		if matched, _ := regexp.MatchString(item.error, err.Error()); !matched {
			t.Fatalf("%d: expect = %v, got = %c", i, item.error, err)
		}
	}
}

var tMarshalTaggedIfaceTest interface{} = &struct{ A string }{"B"}

var tGetterTests = []struct {
	data, tag string
	value     interface{}
}{
	{"_:\n  hi: there\n", "", map[interface{}]interface{}{"hi": "there"}},
	{"_:\n- 1\n- A\n", "", []interface{}{1, "A"}},
	{"_: 10\n", "", 10},
	{"_: null\n", "", nil},
	{"_: !foo BAR!\n", "!foo", "BAR!"},
	{"_: !foo 1\n", "!foo", "1"},
	{"_: !foo '\"1\"'\n", "!foo", "\"1\""},
	{"_: !foo 1.1\n", "!foo", 1.1},
	{"_: !foo 1\n", "!foo", 1},
	{"_: !foo 1\n", "!foo", uint(1)},
	{"_: !foo true\n", "!foo", true},
	{"_: !foo\n- A\n- B\n", "!foo", []string{"A", "B"}},
	{"_: !foo\n  A: B\n", "!foo", map[string]string{"A": "B"}},
	{"_: !foo\n  a: B\n", "!foo", &tMarshalTaggedIfaceTest},
}

func TestMarshalTypeCache(t *testing.T) {
	var data []byte
	var err error
	func() {
		type T struct{ A int }
		data, err = Marshal(&T{})
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()
	func() {
		type T struct{ B int }
		data, err = Marshal(&T{})
		if err != nil {
			t.Fatalf("%v", err)
		}
	}()
	if a, b := string(data), "b: 0\n"; a != b {
		t.Fatalf("expect = %v, got = %v", b, a)
	}
}

type tTypeWithGetter struct {
	tag   string
	value interface{}
}

func (o tTypeWithGetter) GetYAML() (tag string, value interface{}) {
	return o.tag, o.value
}

type typeWithGetterField struct {
	Field tTypeWithGetter "_"
}

func TestMashalWithGetter(t *testing.T) {
	for i, item := range tGetterTests {
		obj := &typeWithGetterField{}
		obj.Field.tag = item.tag
		obj.Field.value = item.value
		data, err := Marshal(obj)
		if err != nil {
			t.Fatalf("%d: %v", i, err)
		}
		if a, b := string(data), string(item.data); a != b {
			t.Fatalf("%d: expect = %v, got = %v%v", i, b, a)
		}
	}
}

func TestUnmarshalWholeDocumentWithGetter(t *testing.T) {
	obj := &tTypeWithGetter{}
	obj.tag = ""
	obj.value = map[string]string{"hello": "world!"}
	data, err := Marshal(obj)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if a, b := string(data), "hello: world!\n"; a != b {
		t.Fatalf("expect = %v, got = %v%v", b, a)
	}
}

func TestSortedOutput(t *testing.T) {
	order := []interface{}{
		false,
		true,
		1,
		uint(1),
		1.0,
		1.1,
		1.2,
		2,
		uint(2),
		2.0,
		2.1,
		"",
		".1",
		".2",
		".a",
		"1",
		"2",
		"a!10",
		"a/2",
		"a/10",
		"a~10",
		"ab/1",
		"b/1",
		"b/01",
		"b/2",
		"b/02",
		"b/3",
		"b/03",
		"b1",
		"b01",
		"b3",
		"c2.10",
		"c10.2",
		"d1",
		"d12",
		"d12a",
	}
	m := make(map[interface{}]int)
	for _, k := range order {
		m[k] = 1
	}
	data, err := Marshal(m)
	if err != nil {
		t.Fatalf("%v", err)
	}

	out := "\n" + string(data)
	last := 0
	for i, k := range order {
		repr := fmt.Sprint(k)
		if s, ok := k.(string); ok {
			if _, err = strconv.ParseFloat(repr, 32); s == "" || err == nil {
				repr = `"` + repr + `"`
			}
		}
		index := strings.Index(out, "\n"+repr+":")
		if index == -1 {
			t.Fatalf("%#v is not in the output: %#v", k, out)
		}
		if index < last {
			t.Fatalf("%#v was generated before %#v: %q", k, order[i-1], out)
		}
		last = index
	}
}

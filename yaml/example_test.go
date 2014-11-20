// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yaml_test

import (
	"fmt"
	"log"

	"github.com/gopkg/encoding/yaml"
)

func ExampleUnmarshal() {
	const data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

	var t struct {
		A string
		B struct {
			C int
			D []int ",flow"
		}
	}
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%v", t)
	// Output:
	// {Easy! {2 [3 4]}}
}

func ExampleUnmarshal_map() {
	const data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

	var m map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("a: %v\n", m["a"])
	fmt.Printf("b/c: %v\n", m["b"].(map[interface{}]interface{})["c"])
	fmt.Printf("b/d: %v\n", m["b"].(map[interface{}]interface{})["d"])
	// Output:
	// a: Easy!
	// b/c: 2
	// b/d: [3 4]
}

func ExampleMarshal() {
	type T struct {
		A string
		B struct {
			C int
			D []int ",flow"
		}
	}

	var t T
	t.A = "Easy!"
	t.B.C = 2
	t.B.D = []int{3, 4}

	d, err := yaml.Marshal(t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s", string(d))
	// Output:
	// a: Easy!
	// b:
	//   c: 2
	//   d: [3, 4]
}

func ExampleMarshal_map() {
	m := map[interface{}]interface{}{
		"a": "Easy!",
		"b": map[interface{}]interface{}{
			"c": 2,
			"d": []int{3, 4},
		},
	}

	d, err := yaml.Marshal(m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s", string(d))
	// Output:
	// a: Easy!
	// b:
	//   c: 2
	//   d:
	//   - 3
	//   - 4
}

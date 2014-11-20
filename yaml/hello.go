// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"github.com/chai2010/yaml"
)

func main() {
	type T struct {
		F int "a,omitempty"
		B int
	}

	d1, _ := yaml.Marshal(&T{B: 2}) // Returns "b: 2\n"
	d2, _ := yaml.Marshal(&T{F: 1}) // Returns "a: 1\nb: 0\n"

	println(string(d1))
	println(string(d2))
}

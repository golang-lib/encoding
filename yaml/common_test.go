// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yaml

type tInlineB struct {
	B        int
	tInlineC `yaml:",inline"`
}

type tInlineC struct {
	C int
}

func tNewInt(x int) *int {
	return &x
}

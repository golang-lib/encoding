// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoding provides encoding extension drivers.
package encoding

import (
	_ "github.com/gopkg/encoding/base58"
	_ "github.com/gopkg/encoding/markdown"
	_ "github.com/gopkg/encoding/toml"
	_ "github.com/gopkg/encoding/yaml"
)

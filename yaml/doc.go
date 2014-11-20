// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package yaml implements YAML support.

Examples:

	func main() {
		var m map[interface{}]interface{}
		err := yaml.Unmarshal([]byte(`... yaml data ...`), &m)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		d, err := yaml.Marshal(m)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}


Source code and other details for the project are available at GitHub:

	https://github.com/chai2010/yaml
*/
package yaml

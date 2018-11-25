// SPDX-License-Identifier: Apache-2.0

package pkg

// Options ...
type Options struct {
	Address    string
	Port       int
	TTL        int
	Listen     bool
	Files      []string
	Token      []byte
	InputPath  string
	OutputPath string
}

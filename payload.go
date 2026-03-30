/*
Copyright © 2026 fakeapate <fakeapate@pm.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	_ "embed"
	"encoding/json"
)

//go:embed data/payloads.json
var payloadData []byte

// payloads holds the traversal strings, URL prefixes, and detection patterns
// loaded from the embedded data/payloads.json file.
type payloads struct {
	Traversals []string          `json:"traversals"`
	Prefixes   []string          `json:"prefixes"`
	Patterns   map[string]string `json:"patterns"`
}

// loadPayloads unmarshals the embedded payloads.json into a [payloads] value.
func loadPayloads() (*payloads, error) {
	var p payloads
	if err := json.Unmarshal(payloadData, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

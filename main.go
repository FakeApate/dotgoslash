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
	"os"
	"strings"

	"github.com/charmbracelet/log"
	flag "github.com/spf13/pflag"
)

func main() {
	var cfg CliConfig

	flag.StringVarP(&cfg.URL, "url", "u", "", "target URL")
	flag.StringVarP(&cfg.Target, "string", "s", "", "substring within --url to replace with payloads")
	flag.StringVarP(&cfg.Cookie, "cookie", "c", "", "Cookie header value")
	flag.IntVarP(&cfg.Depth, "depth", "d", 6, "maximum traversal depth")
	flag.BoolVarP(&cfg.Verbose, "verbose", "v", false, "show all requests")
	flag.Parse()

	if cfg.URL == "" || cfg.Target == "" {
		log.Error("--url and --string are required")
		flag.Usage()
		os.Exit(1)
	}

	if !strings.Contains(cfg.URL, cfg.Target) {
		log.Fatal("target string not found in URL", "string", cfg.Target, "url", cfg.URL)
	}

	if cfg.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	p, err := loadPayloads()
	if err != nil {
		log.Fatal("failed to load payloads", "err", err)
	}

	newScanner(cfg, p).run()
}

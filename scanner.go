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
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fakeapate/mullvad"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
)

// CliConfig holds the runtime parameters parsed from CLI flags.
type CliConfig struct {
	URL     string
	Target  string
	Cookie  string
	Depth   int
	Verbose bool
}

// scanner drives the traversal scan against a single target URL.
type scanner struct {
	cfg      CliConfig
	payloads *payloads
	compiled map[string]*regexp.Regexp
}

// newScanner constructs a [scanner] and pre-compiles all detection regexes.
func newScanner(cfg CliConfig, p *payloads) *scanner {
	compiled := make(map[string]*regexp.Regexp, len(p.Patterns))
	for word, pat := range p.Patterns {
		compiled[word] = regexp.MustCompile(pat)
	}
	return &scanner{cfg: cfg, payloads: p, compiled: compiled}
}

// req is a pre-generated URL and its associated pattern word.
type req struct {
	url  string
	word string
}

// generate builds every (url, word) combination across all depths upfront.
func (s *scanner) generate() []req {
	var reqs []req
	for count := 0; count <= s.cfg.Depth; count++ {
		for _, dvar := range s.payloads.Traversals {
			for _, bvar := range s.payloads.Prefixes {
				for word := range s.payloads.Patterns {
					rewrite := bvar + strings.Repeat(dvar, count) + word
					url := strings.Replace(s.cfg.URL, s.cfg.Target, rewrite, 1)
					reqs = append(reqs, req{url: url, word: word})
				}
			}
		}
	}
	return reqs
}

// run executes the scan from depth 0 up to cfg.Depth.
func (s *scanner) run() {
	c := colly.NewCollector(colly.Async(true))
	mullvadCfg := mullvad.DefaultMullvadConfig()
	c.SetClient(&http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	})

	c.Limit(&colly.LimitRule{Parallelism: 10}) //nolint:errcheck

	if s.cfg.Cookie != "" {
		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Cookie", s.cfg.Cookie)
		})
	}
	extensions.RandomUserAgent(c)

	connected, err := mullvad.IsConnected()
	if err != nil {
		log.Warn("Mullvad connectivity check failed, running without proxies!", "err", err)
	} else if connected {

		mullvad.StartUpdater(mullvadCfg)
		proxies, err := mullvad.SelectProxies(mullvadCfg, 50, mullvad.RelayFilter{
			Weight: func(num int) bool {
				return num <= 99
			},
		})
		if err != nil {
			log.Warn("Failed to select proxies, running without proxies", "err", err)
		} else if p, err := proxy.RoundRobinProxySwitcher(proxies...); err == nil {
			c.SetProxyFunc(p)
		}
	} else {
		log.Warn("Not connected to Mullvad, running without proxies")
	}
	reqs := s.generate()

	c.OnResponse(func(r *colly.Response) {
		word := r.Request.Ctx.Get("word")
		pat, ok := s.compiled[word]
		if !ok {
			return
		}

		rawURL := r.Request.URL.String()
		matches := pat.FindAllString(string(r.Body), -1)

		if len(matches) > 0 {
			fmt.Print("\r\033[K") // clear bar line before logging
			log.Info("match found", "status", r.StatusCode, "url", rawURL, "count", len(matches))
			for i, m := range matches {
				if i >= 6 {
					log.Info("output truncated", "shown", 6, "total", len(matches))
					break
				}
				log.Info("content", "match", m)
			}
		} else {
			log.Debug("response", "status", r.StatusCode, "url", rawURL)
		}

	})

	c.OnError(func(r *colly.Response, err error) {
		log.Debug("request error", "url", r.Request.URL, "err", err)
	})

	log.Info("starting scan", "requests", len(reqs))

	for _, r := range reqs {
		ctx := colly.NewContext()
		ctx.Put("word", r.word)
		c.Request("GET", r.url, nil, ctx, nil)
	}
	c.Wait()

}

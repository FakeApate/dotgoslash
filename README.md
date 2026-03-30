# dot-go-slash

A Go rewrite of [dotdotslash], using [colly] as the HTTP client for proxy rotation support.

Proxy rotation setup is ~~documented~~ in the [mullvad] repo.

```
dotgoslash -url "http://example.com/path/file.pdf" -string "file.pdf" -depth 3 -cookie "session=abc"
```

flags: `-url`/`-u`, `-string`/`-s`, `-cookie`/`-c`, `-depth`/`-d` (default 6), `-verbose`/`-v`

Use responsibly (or don't), I'm not taking any responsibility for the actions a user does with this code.

[dotdotslash]: https://github.com/jcesarstef/dotdotslash
[colly]: https://github.com/gocolly/colly
[mullvad]: https://github.com/fakeapate/mullvad

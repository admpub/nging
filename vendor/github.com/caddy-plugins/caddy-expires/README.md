# caddy-expires

Provides a directive to add expires headers to certain paths or headers

## Usage

Add an expires block to your CaddyFile

```
expires {
    match [a valid regexp] [1y1m1d1h1i1s]
    match_header [header name] [a valid regexp] [1y1m1d1h1i1s]
}
```

Duration can be any combination of y(ear), m(onth), d(ay), h(our), i(minute), s(econd). Parts can be omitted but must
remain in that order

If a path or header matches, an `Expires` header is set (don't forget to version your assets)

## License

MIT

# RequestCache
baa middleware for cache get request based on url and context.

## Options

### Enabled

cache enabled or not, default is `false`.

### Expires

cache expires in seconds, default is `60 * 10`s.

### Headers

set response headers, default is `nil`.

### ContextRelated

default is `false`.

if set `true`, baa context data will be used as request cache key.
data set to baa context by using `c.Set("foo", "bar")` will be json marshaled as part of cache key.

## Requirement

a **cache** DI must be set before use this middleware:

```go
	import (
		"github.com/go-baa/cache"
		_ "github.com/go-baa/cache/redis"
	)

	b.SetDI("cache", cache.New(cache.Options{
		// ...
	}))
```

## Usage

### global
Set as the **last middleware** of baa:
```go
	if baa.Env == baa.PROD {
		// Gzip
		b.Use(gzip.Gzip(gzip.Options{CompressionLevel: 4}))

		// Request Cache
		b.Use(requestcache.Middleware(requestcache.Option{
			Enabled: true,
			Expires: requestcache.DefaultExpires,
		}))
	}
```
Or after request/content processing middlewares.

### with router

```go
	cache := requestcache.Middleware(requestcache.Option{
		Enabled: !b.Debug(),
		Expires: requestcache.DefaultExpires,
	})

	b.Group("/some-prefix", func() {
		// ...
	}, cache)
```

### different options

```go
	cache1 := requestcache.Middleware(requestcache.Option{
		Enabled: !b.Debug(),
		Expires: 60 * 10,
	})

	b.Group("/some-prefix", func() {
		// ...
	}, cache1)

	cache2 := requestcache.Middleware(requestcache.Option{
		Enabled: !b.Debug(),
		Expires: 60 * 30,
	})

	b.Group("/some-prefix-2", func() {
		// ...
	}, cache2)
```
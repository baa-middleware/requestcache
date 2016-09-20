# RequestCache
Cache Get Request

## Options

### Expires

cache expires in seconds, default is 60s * 5.

## Usage

Set as the last middleware of baa:
```
	if baa.Env == baa.PROD {
		// Gzip
		b.Use(gzip.Gzip(gzip.Options{CompressionLevel: 4}))

		// Request Cache
		b.Use(requestcache.Middleware(requestcache.Option{
			Expires: requestcache.DefaultExpires,
		}))
	}
```
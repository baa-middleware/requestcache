package requestcache

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-baa/baa"
	"github.com/go-baa/cache"
)

type response struct {
	Header  map[string][]string
	Body    []byte
	Created time.Time
}

var (
	headerKey = "X-Request-Cache"
)

// Middleware baa middleware func
func Middleware(opt Option) baa.HandlerFunc {
	// get expires option
	var expires = opt.Expires
	if expires == 0 {
		expires = DefaultExpires
	}

	return func(c *baa.Context) {
		// not enabled
		if !opt.Enabled {
			c.Next()
			return
		}

		// only cache get request
		if c.Req.Method != http.MethodGet {
			c.Next()
			return
		}

		// set headers
		if opt.Headers != nil {
			for k, v := range opt.Headers {
				c.Resp.Header().Set(k, v)
			}
		}

		// prepare cache key
		url := c.URL(true)
		if opt.ContextRelated {
			enc, err := json.Marshal(c.Gets())
			if err == nil {
				url = url + ":" + string(enc)
			}
		}
		key := "RequestCache:" + md5Encode(url)

		// read from cache
		cacher := c.DI("cache").(cache.Cacher)
		val := response{}
		if err := cacher.Get(key, &val); err == nil {
			if c.Baa().Debug() {
				c.Baa().Logger().Printf("[RequestCache]: hit [%s]\n", url)
			}

			for k, v := range val.Header {
				for j := range v {
					c.Resp.Header().Set(k, v[j])
				}
			}
			c.Resp.Header().Set(headerKey, "hit")
			c.Resp.Write(val.Body)
			return
		}

		// replace writer
		writer := c.Resp.GetWriter()
		ghostWriter := &ghostWriter{Writer: writer}
		c.Resp.SetWriter(ghostWriter)

		c.Next()

		// restore writer
		c.Resp.SetWriter(writer)

		// non-normal response code
		if c.Resp.Status() != http.StatusOK {
			return
		}

		// skip other content
		switch c.Resp.Header().Get("Content-Type") {
		case baa.ApplicationJSON:
		case baa.ApplicationJSONCharsetUTF8:
		case baa.ApplicationXML:
		case baa.ApplicationXMLCharsetUTF8:
		case baa.TextHTML:
		case baa.TextHTMLCharsetUTF8:
		case baa.TextPlain:
		case baa.TextPlainCharsetUTF8:
		default:
			return
		}

		if c.Baa().Debug() {
			c.Baa().Logger().Printf("[RequestCache]: miss [%s]\n", url)
		}

		// prepare cache content
		val = response{
			Header:  make(map[string][]string),
			Body:    ghostWriter.Content,
			Created: time.Now(),
		}

		// copy header
		header := c.Resp.Header()
		for k, v := range header {
			arr, ok := val.Header[k]
			if !ok {
				arr = make([]string, 0)
			}
			for j := range v {
				arr = append(arr, v[j])
			}
			val.Header[k] = arr
		}

		// cache response
		if err := cacher.Set(key, val, expires); err != nil {
			c.Error(err)
			return
		}

		if c.Baa().Debug() {
			c.Baa().Logger().Printf("[RequestCache]: set [%s]\n", url)
		}
	}
}

func init() {
	gob.Register(response{})
}

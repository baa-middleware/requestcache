package requestcache

const (
	// DefaultExpires default expires in seconds
	DefaultExpires = 60 * 10
)

// Option ...
type Option struct {
	Enabled        bool
	Expires        int64
	Headers        map[string]string
	ContextRelated bool
}

package requestcache

const (
	// DefaultExpires default expires in seconds
	DefaultExpires = 60 * 5
)

// Option ...
type Option struct {
	Expires int64
}
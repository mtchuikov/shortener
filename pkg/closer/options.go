package closer

type Option func(*closer)

func WithMaxConcurrent(max int) Option {
	return func(c *closer) {
		c.maxConcurrent = max
	}
}

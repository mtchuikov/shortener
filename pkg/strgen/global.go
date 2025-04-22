package strgen

var Global *Generator = nil

func InitGlobal(opts ...Option) {
	Global = newGenerator(opts...)
}

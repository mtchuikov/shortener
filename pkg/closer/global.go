package closer

var Global *Closer = nil

func InitGlobal(opts ...Option) {
	Global = New(opts...)
}

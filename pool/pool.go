package pool

// Instance ...
type Instance interface {
	Init()
	Done()
}

// impl is basically a buffered channel
type impl chan struct{}

// Init on single pool
func (p impl) Init() { p <- struct{}{} }

// Done on single pool
func (p impl) Done() { <-p }

// New create a pool instance
func New(size int) Instance {
	if size < 1 {
		size = 1000
	}
	return make(impl, size)
}

// Init on multiple pools
func Init(ps ...Instance) {
	for _, p := range ps {
		p.Init()
	}
}

// Done on multiple pools
func Done(ps ...Instance) {
	for _, p := range ps {
		p.Done()
	}
}

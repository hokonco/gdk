package parallel

// Run ...
func Run(nWait int, buffer int, fns ...func()) {
	var l = len(fns)
	if buffer < 1 {
		buffer = 10_000
	}
	if l > buffer {
		l = buffer
	}
	var p = make(chan struct{}, l)
	for _, fn := range fns {
		go func(fn func()) { fn(); p <- struct{}{} }(fn)
	}
	for i := range fns {
		if nWait == i {
			break
		}
		<-p
	}
}

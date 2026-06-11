package frontier

type Frontier struct {
	queue chan string
}

func New(size int) *Frontier {
	return &Frontier{
		queue: make(chan string, size),
	}
}

func (f *Frontier) Push(url string) {
	f.queue <- url
}

func (f *Frontier) Pop() <-chan string {
	return f.queue
}

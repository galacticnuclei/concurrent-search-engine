package crawler

type Document struct {
	URL     string
	Title   string
	Content string
	Links   []string
}

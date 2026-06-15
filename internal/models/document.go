package models

type Document struct {
	ID      int
	URL     string
	Title   string
	Content string
	Links   []string
}

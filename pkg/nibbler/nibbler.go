package nibbler

// A Nibbler can shorten URLs and also Reverse shortened URLs via the Nibbler API
type Nibbler interface {
	Shorten(url string) (string, error)
	Reverse(id string) (string, error)
}

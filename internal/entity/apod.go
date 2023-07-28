package entity

type Apod struct {
	ID             int
	Date           string
	Explanation    string
	HdUrl          string
	MediaType      string
	ServiceVersion string
	Title          string
	Url            string
	Image          []byte
}

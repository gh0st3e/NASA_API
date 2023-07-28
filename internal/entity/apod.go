package entity

type Apod struct {
	ID             int    `json:"id"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HdUrl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	Url            string `json:"url"`
}

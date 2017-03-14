package models

type JsonBookSearchResult struct {
	Results []Book		`json:"results"`
	MaxCount int		`json:"max-count"`
}

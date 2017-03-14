package models

type JsonBookSearchReq struct {
	Query	string 	`json:"query"`
	Limit	int	`json:"limit"`
}

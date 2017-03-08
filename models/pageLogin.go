package models

type PageLogin struct {
	Title string
	TargetUrl string
	Success	bool
	Error	string
	User	*User
}

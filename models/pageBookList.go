package models

type PageBookList struct {
	Page
	Books 		[]Book
	Error 		string
	PageCount 	int
}

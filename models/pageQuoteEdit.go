package models

type PageQuoteEdit struct {
	Page
	Quote	 	Quotation
	Book		Book
	Error 		string
}
